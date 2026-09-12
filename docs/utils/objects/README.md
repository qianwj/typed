# `github.com/qianwj/typed/utils/objects`

对动态类型 Go 值的工具函数。

## 包导入

```go
import "github.com/qianwj/typed/utils/objects"
```

## `IsNil[T any](value T) bool`

按 Go 的 nil 语义报告 `value` 是否为 nil：

- untyped nil 接口。
- typed nil：指针 / map / slice / channel / function / interface。
- `int` / `string` / 结构体 / 数组 —— 永不为 nil，返回 `false`。

```go
var p *Foo
objects.IsNil(p)             // true
objects.IsNil(nil)           // true
objects.IsNil((*Foo)(nil))   // true
objects.IsNil([]int(nil))    // true
objects.IsNil(0)             // false
```

## `Equals[T any](a, b T) bool`

比 `reflect.DeepEqual` 更人性化的"深度相等"。在三个有意为之的地方不同：

1. **自定义 `Equal` 方法分派。** 如果任一侧的动态类型定义了 `func (T) Equal(T) bool`，会调用该方法，**不走结构递归**。`time.Time` 之类没有声明接口但方法签名匹配的类型也能被识别 —— 反射按"形"找，不按"名"找。如果一侧有 `Equal`、另一侧没有，仍然调用有 `Equal` 的那一侧（两侧都检查，所以 `Equals(a, b) == Equals(b, a)` 在单边实现时仍然成立）。

2. **nil 与空集合的归一化（递归）。** nil 切片与空非 nil 切片在每一层都被视为相等；map 同理。`reflect.DeepEqual` 在这些对上返回 `false`。

3. **typed nil 上的 `Equal` 不会触发 panic。** 若 `Equal` 是指针接收者且任一操作数是 typed nil 指针，两边都 nil 返回 `true`，否则 `false`，**绝不**调用 nil 接收者的方法。

其它方面沿用 `reflect.DeepEqual`：解引用指针、字段逐个比（含未导出字段）、数组逐元素比、循环结构不爆栈。

### 启用自定义 `Equal`

只要你的类型有签名严格为 `func (T) Equal(T) bool` 的方法即可，**不需要**实现 `objects.Equaler`（该接口已保留用于文档/类型提示，但 `Equals` 不会要求）。签名不匹配（例如 `Equal(Interface) bool`）时 `Equals` 不会采用，会回退到结构递归。

```go
type Vec2 struct{ X, Y int }
func (a Vec2) Equal(b Vec2) bool { return a.X == b.X && a.Y == b.Y }

objects.Equals(Vec2{1, 2}, Vec2{1, 2}) // true（走 Equal）
objects.Equals([]int(nil), []int{})    // true（nil/空归一化）
```
