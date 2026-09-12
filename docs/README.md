# Typed · Documentation

> Per-package API reference and examples for the **Typed** toolkit.
>
> [中文索引](./README-cn.md) · [Back to repo README](../README.md)

Each file follows a five-section shape — **Import → Why → API → Example → See also** — so it can be read on its own. Every doc is also mirrored in Chinese as `README-cn.md` next to it.

## Contents

### Collections and synchronous streams

- 📚 [`collections`](./collections/README.md) — `Stack` / `Queue` / `Deque` / `ArrayList` / `LinkedList` / `HashMap` / `HashSet` / `Stream`, plus `collections.Range`.
- 📡 [`reactivex`](./reactivex/README.md) — asynchronous event streams: `Observable` / `Publisher` / `Subscriber` / `Subject`, backpressure, sources and operators.

### Control flow

- 🔁 [`control`](./control/README.md) — `Repeat` / `RepeatE` and `control/match` pattern matching (`Pattern[T]` plus chained `Case` / `Type`).

### Abstractions

- 🟢 [`utils/option`](./option/README.md) — `Optional[T]`, a value that may be absent.
- 🟡 [`utils/result`](./result/README.md) — `Result[T]`, success / failure, and bridges to and from `(T, error)`.
- 🟰 [`utils/objects`](./utils/objects/README.md) — `IsNil` / `Equals`, with `Equal(T) bool` dispatch and nil-vs-empty normalisation.
- 📦 [`utils/json`](./utils/json/README.md) — `Encode` / `Decode`, a `Result`-style codec on top of `encoding/json/v2`.

## Reading paths

- 🆕 **New to Typed** — start with [`collections`](./collections/README.md) to see the immutable-transform style of `ArrayList` / `LinkedList` / `Stream`, then read [`utils/option`](./option/README.md) and [`utils/result`](./result/README.md) for the absent / failure story.
- ⚡ **Asynchronous or multi-subscriber behavior** — jump to [`reactivex`](./reactivex/README.md).
- 🧭 **Value- or type-based dispatch** — jump to [`control`](./control/README.md).
