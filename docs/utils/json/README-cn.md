# `github.com/qianwj/typed/utils/json`

基于 [`encoding/json/v2`](https://pkg.go.dev/encoding/json/v2) 的 `Result` 风格编解码。

## 包导入

```go
import (
    "github.com/qianwj/typed/utils/json"
    "github.com/qianwj/typed/utils/result"
    jsonopts "github.com/go-json-experiment/json/options" // 可选
)
```

## `Encode[T any](t T, opts ...json.Options) result.Result[[]byte]`

把值 `t` 序列化为 JSON 字节切片并以 `Result[[]byte]` 返回：

- 成功路径：携带 JSON 字节。
- 失败路径：携带 `encoding/json/v2` 的错误（**未包装**）。需要区分 JSON 错误时用 `errors.Is` / `errors.As` 对照 v2 错误类型。

`Encode` 内部用 `result.Wrap` 适配 `json.Marshal` 的 `([]byte, error)`，所以一次调用就拿到 `Result` —— 调用点不需要 `if err != nil`。

### 选项

`opts` 是 `encoding/json/v2` 的可变 `json.Options`（`jsonopts.Options`）。每个选项是一个 property setter，后面的覆盖前面的。常见用途：

```go
data := json.Encode(value, json.Deterministic(true))
```

v2 选项不是 v1 `json.Marshal` 的 tag（如 `MarshalJSON` / `UnmarshalJSON`）。v2 的选项是运行时配置；tag 与 `Marshaler` 方法依然是控制产物内容的标准方式。

### 与 `Decode` 的往返

`Encode` 与 `Decode` 互逆：相同目标类型往返一次后，**按 Go `==` 比较解码字段**得到相等值（除 channel / function / complex 这些 JSON 不可表达的字段外）。这是项目其它部分依赖的契约；`encoder_test.go` 显式覆盖了它。

### 内存

`Encode` 会**为编码结果分配全新的 `[]byte`**，调用方拥有这段切片并可在消费 `Result` 后立即复用 / 池化 / 释放。输入值 `t` 在 `Encode` 返回后即可修改 —— `json.Marshal` 不会保留它。

### HTML 转义

与 v1 不同，v2 默认**不做** HTML 转义。值中含 `<` / `>` / `&` 时会被原样写出。需要 HTML 转义时请单独跑一遍转义步骤，或显式给 v1 风格选项。

## `Decode[T any](data []byte) result.Result[T]`

把 `data` 解码到类型 `T` 并以 `Result[T]` 返回。

- 成功路径：解码后的 `T` 值。
- 失败路径：底层 `encoding/json/v2` 错误（**未包装**）。

```go
r := json.Decode[Config](raw)
if r.IsFailure() {
    return r.Unwrap() // (Config 零值, err)
}
cfg := r.Value()
use(cfg)
```

## 例子

```go
import (
    "github.com/qianwj/typed/utils/json"
    "github.com/qianwj/typed/utils/result"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

encoded := json.Encode(User{ID: 1, Name: "Ada"})
encoded.IfPresent(func(b []byte) { fmt.Println(string(b)) })

decoded := json.Decode[User](encoded.OrElseGet(func() []byte { return nil }))
name := decoded.Map(func(u User) string { return u.Name }).
    OrElse("<unknown>")
```

## 与其他包的关系

- 错误与成功流用 [`utils/result`](../../result/README-cn.md)；本包只通过 `Result[T]` 与外界交互。
