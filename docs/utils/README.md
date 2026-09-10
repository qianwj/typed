# typed/utils 设计说明

## 模块定位

`typed/utils` 是 `typed` 项目下独立的工具模块,提供两个泛型值类型:

- `Optional[T]` —— 表示"可能不存在的值"
- `Result[T]` —— 表示"可能失败的结果"

它不依赖 `typed/collections`,也不被 `typed/collections` 直接依赖。两个模块可以通过类型本身自然组合(`Result.Optional()` 返回 `Optional[T]`),但保持独立演进。

预计独立发布为:

```text
github.com/qianwj/typed/utils
```

## 为什么需要这两个类型

Go 标准库处理"值不存在"的常用方式是 `(T, bool)`,处理"可能失败"的方式是 `(T, error)`。它们在简单调用点足够清晰,但当一段流程需要多步组合时,就会变成反复的 `if err != nil` / `if ok { ... }` 检查,处理流程被这些守卫打散。

`Optional` 和 `Result` 提供三个价值:

1. **自描述的类型**:调用方从签名就能看出"这一步可能没有值"或"这一步可能失败",不需要去看文档或注释。
2. **可组合的链式 API**:`Map` / `FlatMap` / `Filter` / `OrElse` 让流程可以左到右写下去,不再需要为每一步写守卫。
3. **明确的语义边界**:`Optional` 只关心"有没有",`Result` 同时关心"是哪个 error",避免一个类型承担两个职责。

## Optional[T]

### 状态表示

`Optional` 内部用一个 `present bool` 标记加上 `value T` 字段:

```go
type Optional[T any] struct {
    value   T
    present bool
}
```

这跟 Java 的 `java.util.Optional` 一致。**不能**靠"`value == nil` 判断空"——那样值类型(`int`、`string`、`struct{}`)在编译期就会失败,只能用于指针/接口/slice/map/chan/func。

测试 `TestOptionalOf/zero_value_is_still_present` 显式覆盖了 `Of(0)` 仍然是 present 这一关键不变量,这是这套设计能成立的根本原因。

### 构造器

| 构造器 | 适用场景 |
| --- | --- |
| `Empty[T]()` | 明确表示"无值",适合 `int`、`struct` 等值类型 |
| `Of(value)` | 已知有值,包括零值 |
| `OfNullable(value)` | 指针/slice/map 等可空类型;nil 视为空 |

`OfNullable` 内部用 `isNil` 反射判断。`isNil` 同时处理 typed nil 和 untyped nil interface(`var i any`)两种情况,因为后者会让 `reflect.ValueOf` 返回 `Kind() == Invalid` 的零 Value,被任何显式的 Kind switch 漏掉。

### 关键方法

| 方法 | 行为 |
| --- | --- |
| `IsPresent` / `IsEmpty` | 查询状态 |
| `Get` | 取值,空时 panic(配套 `IsPresent` 用) |
| `OrElse(default)` | 空时返回默认值(立即求值) |
| `OrElseGet(f)` | 空时调用 `f`(惰性求值) |
| `OrElseThrow(msg)` | 空时返回 `(zero, errors.New(msg))` |
| `IfPresent(f)` | 有值时执行副作用 |
| `IfPresentOrElse(p, a)` | 两个分支,恰好跑一个 |
| `Filter(p)` | 有值且满足谓词时保留,否则变空 |
| `Map[R](f)` | 类型变换 |
| `FlatMap[R](f)` | 嵌套 Optional 展开 |

`Map` 和 `FlatMap` 是 **Go 1.27 泛型方法**——它们在自己的方法签名上声明 `[R any]`,而 `Optional` 作为具体泛型类型可以承载。接口方法至今不能声明类型参数,这也是 `Optional` 选择 struct 而不是 interface 的原因,跟 `ArrayList[T]` 的设计一致。

### 命名取舍

- `OrElseThrow` 名字来自 Java 的 `Optional.orElseThrow`,但行为返回 `(T, error)` 而不是 panic。Go 的惯例是显式 error,所以选了 `(T, error)` 形态;文档里明确说明这个名字不暗示 panic。
- `Value` 作为方法名跟 `Value` 包级构造器同名在 Go 里合法(方法接收者 vs 包级名字属于不同命名空间),但读起来仍然别扭。`ArrayListOf` 风格的 `Of` / `Empty` / `OfNullable` 才是真正的入口。

## Result[T]

### 状态表示

```go
type Result[T any] struct {
    value T
    err   error
}
```

`err == nil` 即成功。`Failure[T](nil)` 直接 panic——一个带 nil error 的 Result 应当用 `Success` 构造,这是构造期不变量。

### 命名取舍:`Success` / `Failure` vs `Ok` / `Err`

构造器最终选了 `Success` / `Failure` 而不是更短的 `Ok` / `Err`,理由:

- `Ok` 是缩写,在英语语境里读起来像"okay",不够正式;`Success` 完整拼写,在文档和 PR diff 里都更清楚。
- `Err` 不拼读任何单词,且跟方法 `r.Error()` 视觉撞车——`utils.Err[int](e)` 和 `r.Error()` 在同一份代码里交替出现会增加阅读成本。
- `Failure` 比 `Err` 长,但跟 `Success` 对称,predicate `IsSuccess` / `IsFailure` 也保持一致,整组 API 形成一个命名家族。

> 跟 `Ok` 一起的内置方法 `MapErr` / `GetOrElseErr` 也对应改成 `MapError` / `GetOrElseFailure`,理由同上,避免一个类型里"Err"既是构造器又是方法名片段。

### `OrElse*` 命名家族:与 `Optional` 对齐

值提取器族,失败时用 fallback 直接返回 `T`:

- `OrElse(defaultValue T) T` — 固定默认值(立即求值)
- `OrElseGet(f func() T) T` — 惰性 fallback,看不到 error
- `Recover(f func(error) T) T` — 错误感知的 fallback,拿到 error 才能根据错误类型返回对应的 T

`OrElse*` 这一族名字公认意味着"取一个值出来",跟 Java 的 `java.util.Optional.orElse` / `orElseGet` 一致;`Recover` 是 `Result` 独有的,因为 `Optional` 没有 error 通道。

`Recover` 的典型用法是按错误类型分流:

```go
v := loadProfile(id).Recover(func(err error) Profile {
    if errors.Is(err, ErrNotFound) {
        return Profile{}              // 缺失给零值
    }
    return Profile{Source: "cache"}  // 其他错误兜底走缓存
})
```

| 返回值 | 提取器 |
| --- | --- |
| `T` | `OrElse(d)` / `OrElseGet(f)` / `Recover(f)` |
| `(T, error)` | `Unwrap()` |

如果 fallback 本身可能失败,正确的写法是 `Unwrap()` 之后用 Go 原生 `(T, error)` 链式处理,而不是再造一个 combinator:

```go
v, err := loadProfile(id).FlatMap(validate).Unwrap()
if err != nil {
    v, err = loadFromCache(id) // 任意 (T, error) 函数都可作 fallback
}
return v, err
```

这条路径比再造 combinator 通用——任何 Go 函数都能直接接上,不需要先包成 `Result`,且用 Go 习惯的 `if err != nil` 写后续分支,代码更线性。

#### 三个值提取器怎么选

| 场景 | 用 |
| --- | --- |
| 失败总是给同一个默认值 | `OrElse(d)` |
| 失败给一个可能昂贵的默认值(构造有副作用) | `OrElseGet(f)` |
| 失败后根据错误类型返回不同值(`ErrNotFound` 给零值,其他给缓存) | `Recover(f)` |
| 失败后用 Go 原生错误处理(可能再次失败) | `Unwrap()` |

### 为什么不是 `Result[T, E]`

原本签名是 `Result[T any, E error]`,但 `E` 被约束为 `error` 接口,内部存的也是 `error`,`E` 不携带任何额外类型信息,`MapError[F]` 也写不出有意义的形态。`Result[T]` 简化了类型,又兼容所有 `error` 实现,不损失表达能力。

### 关键方法

| 方法 | 行为 |
| --- | --- |
| `Success(v)` / `Failure(e)` | 构造 |
| `IsSuccess` / `IsFailure` | 查询状态 |
| `Value` | 取值,失败时 panic(panic 值是原始 error) |
| `Error` | 取错误 |
| `Unwrap` | 返回 `(T, error)`,Result 转 Go 标准错误处理的桥 |
| `Optional` | 成功转 present Optional,失败转 empty Optional(丢弃 error) |
| `OrElse(d)` | 失败时返回默认值 `d`(立即求值) |
| `OrElseGet(f)` | 失败时调用 `f`(惰性求值) |
| `Recover(f)` | 失败时调用 `f(err)`,fallback 可基于 error 类型返回不同 T |
| `Map[R](f)` | 成功时变换,失败传播 |
| `FlatMap[R](f)` | 嵌套 Result 展开 |
| `MapError(f)` | 失败时用 `f(err)` 包装 error,常用于在链尾加上下文标签 |

### panic 行为

`Value` 在失败时 `panic(r.err)`,recover 拿到的就是原始 `error`,不会因为 `fmt.Sprintf` 之类丢失类型信息。`TestResultValuePanic` 显式断言了这一点——上层 recover 后还想 `errors.As` 拿具体错误类型,这个契约不能破坏。

## 与标准库的关系

- `(T, bool)` 和 `(T, error)` **仍然是首选**。这两个类型是为多步组合而存在的,不是来替代标准返回形态。
- `Result` 的成功路径用 `T`,失败路径用 `error`,所以从 `Result[T]` 转回 `(T, error)` 的标准做法是 `r.Unwrap()`——它不 panic,把错误交给调用方处理。
- `Optional` 不强制任何 sentinel,Java Optional 用 `null` 表达 absent 的反例一直被人诟病,这里用 `present bool` 避免了类似问题。

## 与 typed 其他模块的关系

- `typed/collections` 不依赖 `utils`,`utils` 也不依赖 `collections`。
- `Result.Optional()` 提供 `Result[T] -> Optional[T]` 的单向桥。
- 后续如果 `Stream[T]` 引入错误传播(`MapE`、`FilterE`),内部可以选用 `Result` 作为错误承载——这是 `Result` 优先做对的一个理由。

## 未来可能的变化

- `Optional` 增加 `Or(other)`(返回另一个 Optional) 和 `Stream() iter.Seq[T]`(只有 present 时产出一个元素)。
- `Result` 增加 `Tap(f func(T))` 和 `TapError(f func(error))` 形式的副作用,对应 Java `peek` 的语义。
- 如果某天业务方真的需要 `Result[T, E]`(区分多种 error 类型),会作为 `Result` 旁边的姐妹类型提供,而不是替换。

## 许可

本模块使用 [MIT License](../../LICENSE)。
