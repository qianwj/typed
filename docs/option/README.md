# typed/utils/option 设计说明

## 模块定位

`typed/utils/option` 是 `typed` 项目下独立的泛型值类型模块,提供 `Optional[T]`,表示"可能不存在的值"。

预计独立发布为:

```text
github.com/qianwj/typed/utils/option
```

`Optional` 与同模块的 [`Result`](../result/README.md) 形成互补: `Optional` 只关心"有没有",`Result` 关心"有没有" + "是哪个 error"。

## 为什么需要 `Optional`

Go 标准库处理"值不存在"的常用方式是 `(T, bool)`。在简单调用点足够清晰,但当一段流程需要多步组合时,就会变成反复的 `if ok { ... }` 检查,处理流程被这些守卫打散。

`Optional` 提供三个价值:

1. **自描述的类型**:调用方从签名就能看出"这一步可能没有值",不需要去看文档或注释。
2. **可组合的链式 API**:`Map` / `FlatMap` / `Filter` / `OrElse` 让流程可以左到右写下去,不再需要为每一步写守卫。
3. **明确区分"零值"和"缺失"**:`Optional` 用 `present bool` 标记,不会让存了零值和缺失状态混淆。

## 状态表示

`Optional` 内部用一个 `present bool` 标记加上 `value T` 字段:

```go
type Optional[T any] struct {
    value   T
    present bool
}
```

这跟 Java 的 `java.util.Optional` 一致。**不能**靠"`value == nil` 判断空"——那样值类型(`int`、`string`、`struct{}`)在编译期就会失败,只能用于指针/接口/slice/map/chan/func。

测试 `TestOptionalOf/zero_value_is_still_present` 显式覆盖了 `Of(0)` 仍然是 present 这一关键不变量,这是这套设计能成立的根本原因。

## 构造器

| 构造器 | 适用场景 |
| --- | --- |
| `Empty[T]()` | 明确表示"无值",适合 `int`、`struct` 等值类型 |
| `Of(value)` | 已知有值,包括零值 |
| `OfNullable(value)` | 指针/slice/map 等可空类型;nil 视为空 |

`OfNullable` 内部用 `isNil` 反射判断。`isNil` 同时处理 typed nil 和 untyped nil interface(`var i any`)两种情况,因为后者会让 `reflect.ValueOf` 返回 `Kind() == Invalid` 的零 Value,被任何显式的 Kind switch 漏掉。

## 关键方法

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

## 命名取舍

- `OrElseThrow` 名字来自 Java 的 `Optional.orElseThrow`,但行为返回 `(T, error)` 而不是 panic。Go 的惯例是显式 error,所以选了 `(T, error)` 形态;文档里明确说明这个名字不暗示 panic。
- 构造器叫 `Of` / `Empty` / `OfNullable`,而不是 `Value` 这类容易跟方法名撞车的词。

## 与标准库的关系

`(T, bool)` 仍然是简单调用点的首选。`Optional` 是为多步组合而存在的,不是来替代标准返回形态。

## 与 typed 其他模块的关系

- `typed/collections` 的 `ArrayList.First` / `ArrayList.Last` / `ArrayList.Find`、`LinkedList.First` / `LinkedList.Last` / `LinkedList.Find`、`Stream.First` / `Stream.Last` / `Stream.Find`、`HashSet.Find` 全部返回 `option.Optional[T]`,统一了"找不到"和"没有"两种语义。
- `typed/utils/result` 的 `Result.Optional()` 提供 `Result[T] -> Optional[T]` 的单向桥。

## 未来可能的变化

- `Optional` 增加 `Or(other)`(返回另一个 Optional) 和 `Stream() iter.Seq[T]`(只有 present 时产出一个元素)。

## 许可

本模块使用 [MIT License](../../../LICENSE)。
