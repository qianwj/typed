# Typed · Documentation

> Per-package API reference and examples for the **Typed** toolkit.
>
> [中文索引](./README-cn.md) · [Back to repo README](../README.md)

Each file follows a five-section shape — **Import → Why → API → Example → See also** — so it can be read on its own. Every doc is also mirrored in Chinese as `README-cn.md` next to it.

## Contents

### Collections and synchronous streams

- 📚 [`collections`](./collections/README.md) — `Stack` / `Queue` / `Deque` / `ArrayList` / `LinkedList` / `HashMap` / `HashSet` / `Stream`, plus `collections.Range`.

### Reactive streams and asynchronous results

- 📡 [`reactivex`](./reactivex/README.md) — `Flowable`, `Single`, `Maybe`, and `Subject`; demand, backpressure, shared cached results and callback composition.
- [Type conversions](./reactivex/README.md#type-conversions) — `ToFlowable`, `FirstElement`, `FirstOrError`, and finite-stream reduction.

### Control flow

- 🔁 [`control`](./control/README.md) — `If` / `IfGet` conditional values, `Repeat` / `RepeatE` loops, and `control/match` pattern matching (`Pattern[T]` plus chained `Case` / `Type`).

### Concurrency primitives

- 🧵 [`concurrency`](./concurrency/README.md) — blocking queues, `Group`, `Semaphore`, and typed temporary object reuse with `Pool[T]`.

### Abstractions

- 🟢 [`adt`](./adt/README.md) — `Option[T]`, a value that may be absent.
- 🟡 [`adt`](./adt/README.md) — `Result[T]`, success / failure, and bridges to and from `(T, error)`.
- 🟰 [`utils/objects`](./utils/objects/README.md) — `IsNil` / `Equals`, with `Equal(T) bool` dispatch and nil-vs-empty normalisation.
- 📦 [`utils/json`](./utils/json/README.md) — `Encode` / `Decode`, a `Result`-style codec on top of `encoding/json/v2`.

## Reading paths

- 🆕 **New to Typed** — start with [`collections`](./collections/README.md) to see the immutable-transform style of `ArrayList` / `LinkedList` / `Stream`, then read [`adt`](./adt/README.md) and [`adt`](./adt/README.md) for the absent / failure story.
- ⚡ **Asynchronous results or streams** — use [`Single`](./reactivex/README.md#singlet) for one value or an error, [`Maybe`](./reactivex/README.md#maybet) for an optional value, [`Flowable`](./reactivex/README.md#flowablet) for multiple values with demand, and [`Subject`](./reactivex/README.md#subjectt) for multicast.
- 🧭 **Value- or type-based dispatch** — jump to [`control`](./control/README.md).
- 🚦 **Bounded producer / consumer queue with backpressure** — jump to [`concurrency`](./concurrency/README.md).
