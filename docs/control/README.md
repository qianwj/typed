# typed/control 使用说明

`typed/control` 提供 Go 泛型下的控制流工具。当前先实现 `match` 包，用来表达模式匹配。

模块路径：

```text
github.com/qianwj/typed/control
```

该模块使用 Go 1.27 的泛型方法，因此编译器和 GoLand SDK 需要使用 Go 1.27 或更高版本。

导入：

```go
import "github.com/qianwj/typed/control/match"
```

## 设计目标

Go 已经有清晰的 `if` 和 `switch`，简单分支应继续使用语言内置控制流。`match` 的目标是处理这些场景：

- 多个分支都要返回同一种结果类型。
- 分支需要组合相等、集合、范围、谓词等匹配条件。
- 输入是 `any`，需要按动态类型分派。
- 调用点希望把匹配逻辑写成从上到下的表达式，并显式处理未匹配分支。

这不是语言级模式匹配，不能做编译期穷尽性检查，也不会做结构体字段解构。Go 里这两件事需要语言支持才可靠。

## 值匹配

值匹配从 `match.Match(value)` 开始，也可以使用等价的 `match.Value(value)`。每个分支返回同一种结果类型 `R`，匹配按声明顺序执行，第一条命中的分支获胜。

```go
grade := match.Value(score).
    Case(match.Between(90, 100), func(int) string {
        return "A"
    }).
    Case(match.Between(80, 89), func(int) string {
        return "B"
    }).
    When(func(n int) bool {
        return n >= 60
    }, func(int) string {
        return "C"
    }).
    Default(func(int) string {
        return "F"
    })
```

值匹配支持的基础模式：

| 模式 | 约束 | 行为 |
| --- | --- | --- |
| `Predicate(p)` | `T any` | 用自定义谓词判断 |
| `Any()` | `T any` | 总是匹配 |
| `Eq(value)` | `T any` | 使用 `objects.Equals` 深度比较，并支持类型自定义的 `Equal` |
| `In(values...)` | `T comparable` | 判断值是否在集合中 |
| `Between(min, max)` | `cmp.Ordered` | 判断是否在闭区间内 |
| `Not(pattern)` | `T any` | 取反 |
| `And(patterns...)` | `T any` | 所有模式都匹配 |
| `Or(patterns...)` | `T any` | 任一模式匹配 |

`Eq` 接受任意类型，使用 `objects.Equals` 判断，因此可比较普通值、struct，以及包含 slice/map 的 struct；如果类型提供 `Equal` 方法，也会使用类型自己的相等语义。需要只对当前分支生效的自定义规则时使用 `When` 或 `Predicate`。`In` 仍要求 `comparable`，因为它使用哈希集合判断成员关系：

```go
size := match.Value([]int{1, 2, 3}).
    When(func(values []int) bool {
        return len(values) > 0
    }, func(values []int) int {
        return len(values)
    }).
    OrElse(0)
```

## 类型匹配

类型匹配从 `match.Type(value)` 开始，输入类型是 `any`。`Case` 的目标类型从 handler 参数推导：

```go
message := match.Type(value).
    Case(func(s string) string {
        return "string: " + s
    }).
    Case(func(n int) string {
        return fmt.Sprintf("int: %d", n)
    }).
    Default(func(value any) string {
        return fmt.Sprintf("unknown: %T", value)
    })
```

如果 `Case` 的参数是接口类型，动态值实现该接口时会命中：

```go
kind := match.Type(err).
    CaseWhen(func(e interface{ Timeout() bool }) bool {
        return e.Timeout()
    }, func(interface{ Timeout() bool }) string {
        return "timeout"
    }).
    Default(func(any) string {
        return "other"
    })
```

`Nil` 分支可以匹配 untyped nil interface，也可以匹配 typed nil pointer、map、slice、channel、function 或 interface：

```go
kind := match.Type(value).
    Nil(func() string {
        return "nil"
    }).
    Case(func(s string) string {
        return "string"
    }).
    Default(func(any) string {
        return "other"
    })
```

typed nil 同时带有动态类型，因此分支顺序会影响结果。如果 `Case(func(*User) R)` 放在 `Nil` 前面，typed nil `*User` 会先命中类型分支；如果 `Nil` 放在前面，会先命中 nil 分支。

## 终止方式

匹配链有以下常用终止方式：

| API | 行为 |
| --- | --- |
| `Default(func(T) R)` / `Default(func(any) R)` | 未匹配时执行默认分支并返回 `R` |
| `OrElse(value)` | 未匹配时返回给定默认值 |
| `OrElseGet(func(...) R)` | 未匹配时惰性计算默认值 |
| `Unwrap() (R, bool)` | 返回匹配结果和是否命中 |
| `Matched()` | 只查询是否已有分支命中 |

`Default`、`OrElseGet` 和分支 handler 都是惰性的：前面已有分支命中时，后续 pattern、guard、handler 和 default 都不会执行。

## 边界

`match` 是表达式级控制流，不替代语言内置 `switch`。简单常量分支仍然优先写 `switch`，需要把匹配作为值传递、需要组合模式、或者需要按动态类型分派时再使用 `match`。
