# `typed/utils/json` — Typed

`Result`-style JSON codec backed by [`encoding/json/v2`](https://pkg.go.dev/encoding/json/v2).

> Part of the **Typed** toolkit. Looking for the Chinese version? See [README-cn.md](./README-cn.md).

## Contents

- [Import](#import)
- [`Encode[T]`](#encodet-t-t-opts-jsonoptions-resultresultbyte)
- [`Decode[T]`](#decodet-data-byte-resultresultt)
- [Example](#example)
- [See also](#see-also)

## Import

```go
import (
    "github.com/qianwj/typed/utils/json"
    "github.com/qianwj/typed/utils/result"
    jsonopts "github.com/go-json-experiment/json/options" // optional
)
```

## `Encode[T any](t T, opts ...json.Options) result.Result[[]byte]`

Serialises a value of type `T` into JSON bytes and returns it as a `Result[[]byte]`:

- **Success:** the JSON byte slice.
- **Failure:** the underlying `encoding/json/v2` error, **unwrapped**. To distinguish JSON errors from other errors, use `errors.Is` / `errors.As` against the v2 error types.

`Encode` is implemented in terms of `result.Wrap` — the `([]byte, error)` pair from `json.Marshal` is wrapped into a `Result` in one call, no `if err != nil` ladder at the call site.

### Options

The `opts` parameter is a variadic of `encoding/json/v2` options (`jsonopts.Options`). Each option is a property setter; later options override earlier ones. Common cases:

```go
data := json.Encode(value, json.Deterministic(true))
```

v2 options are not v1 `json.Marshal` tags like `MarshalJSON` / `UnmarshalJSON`. v2's options are runtime configuration; tags and `Marshaler` methods remain the standard way to influence what is produced.

### Round-trip with `Decode`

`Encode` and `Decode` are inverses. A value encoded by `Encode` and then decoded back into the same target type with `Decode` produces an equal value (by Go's `==` on the decoded fields, modulo the JSON marshalling lossiness for channels, functions, and complex numbers). This is the contract the rest of the project relies on; the tests in `encoder_test.go` cover it explicitly.

### Memory

`Encode` allocates a **fresh `[]byte`** for the encoded value. The caller owns the returned slice and may reuse, pool, or free it as soon as the `Result` is consumed. The input value `t` is read by `json.Marshal` but not retained; it can be mutated as soon as `Encode` returns.

### HTML escaping

Unlike v1, `encoding/json/v2` does **not** escape HTML by default. A value containing `<`, `>`, or `&` is encoded verbatim. If HTML escaping is required, run the bytes through a separate escaping step or use the v1-style options explicitly.

## `Decode[T any](data []byte) result.Result[T]`

Decodes `data` into a value of type `T` and returns it as a `Result[T]`.

- **Success:** the decoded `T` value.
- **Failure:** the underlying `encoding/json/v2` error, **unwrapped**.

```go
r := json.Decode[Config](raw)
if r.IsFailure() {
    return r.Unwrap() // (zero Config, err)
}
cfg := r.Value()
use(cfg)
```

## Example

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

## See also

- The error / success flow uses [`utils/result`](../../result/README.md); this package only interacts with the outside world through `Result[T]`.
