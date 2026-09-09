# typed

`typed` is a type-safe collection toolkit built with Go generics.

Its first goal is to provide fluent, Java/JavaScript-inspired collection types such as `ArrayList[T]` and `HashMap[K, V]`. A lazy `Stream[T]` layer can be used when a pipeline should be evaluated on demand.

> The project is currently in an early design stage. `collections/stream` is still a placeholder API, so the fluent examples below describe the intended direction and may not compile yet.

## Goal

Go's `for` loop is clear and should remain the preferred choice for simple logic. However, when a collection goes through several transformations, nested functions or repeated temporary slices can make the processing flow harder to read:

```go
result := Collect(
    Map(
        Filter(users, func(u User) bool { return u.Age >= 18 }),
        func(u User) string { return u.Name },
    ),
)
```

`typed` aims to support a left-to-right collection style:

```go
result := typed.ArrayList[User](users).
    Filter(func(u User) bool {
        return u.Age >= 18
    }).
    Map(func(u User) string {
        return u.Name
    })
```

For lazy processing, the same collection can become a stream explicitly:

```go
result := typed.ArrayList[User](users).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

## Design Principles

- **Type safety**: Use Go generics and avoid `any`, reflection, and runtime type assertions wherever possible.
- **Collection first**: `ArrayList[T]` and `HashMap[K, V]` are the primary user-facing types for ordinary collection work.
- **Fluent calls**: Collection and stream operations return typed values, making the data flow readable from left to right.
- **Optional laziness**: Collection operations are straightforward and eager; `Stream[T]` is available for lazy pipelines and early termination.
- **Composability**: Collections, iterators, and other Streams can be combined into new data sources.
- **Early termination**: Operations such as `First`, `Any`, and `Take` should avoid processing unnecessary values.
- **Go style**: Do not copy every Java Stream semantic blindly; simple logic should remain easy to write with `for range`.

## API Draft

### Collections

The planned collection types are concrete generic types rather than interfaces:

```go
type ArrayList[T any] []T
type HashMap[K comparable, V any] map[K]V
```

This allows type-changing generic methods such as:

```go
names := ArrayList[User](users).
    Filter(isAdult).
    Map(func(u User) string { return u.Name })
```

`ArrayList[T]` is intended for slice-like ordered data. `HashMap[K, V]` is intended for key-value operations such as `Filter`, `MapValues`, `Keys`, `Values`, and `ToSlice`. The exact naming and eager/lazy boundary are still under design.

### Lazy Streams

`Stream[T]` is an optional lazy layer backed by Go's iterator conventions. It is useful for large, one-shot, channel-backed, or potentially infinite sources:

```go
profiles := ArrayList[User](users).
    Stream().
    Filter(isAdult).
    Map(toProfile).
    Collect()
```

The stream API is not intended to replace the collection types or ordinary `for range` loops.

Planned operations fall into three broad groups.

### Intermediate Operations

These operations generally return another collection or `Stream`, depending on the receiver:

- `Filter`
- `Map`
- `FlatMap`
- `Distinct`
- `Take` / `Drop`
- `Skip` / `Limit`
- `Concat`
- `Peek`

### Terminal Operations

These operations consume the Stream:

- `Collect`
- `Count`
- `First` / `Last`
- `Any` / `All` / `None`
- `Find`
- `Reduce`
- `ToMap`
- `GroupBy`
- `ForEach`

### Sorting and Aggregation

- `Sort` / `SortBy`
- `Min` / `Max`
- `Sum` / `Average`
- `GroupBy`
- `PartitionBy`
- `Join`

The concrete API will be refined around Go error handling, generic method support, and lazy iterator semantics.

## Relationship to go-linq

[`go-linq`](https://github.com/ahmetb/go-linq) is an important reference project and already provides a complete typed LINQ implementation for Go 1.27. Its central abstraction is `Query[T]`, a lazy query value with a broad LINQ-style operator set such as `Where`, `Select`, `GroupBy`, `Join`, and `Aggregate`.

`typed` intentionally takes a different top-level approach:

| Concern | `go-linq` | `typed` direction |
| --- | --- | --- |
| Primary abstraction | Lazy `Query[T]` | Concrete `ArrayList[T]` and `HashMap[K, V]` first |
| Naming | LINQ-oriented: `Where`, `Select` | Collection-oriented: `Filter`, `Map` |
| Evaluation | Lazy by default | Eager collections, explicit lazy `Stream[T]` |
| JavaScript-style arrays | Adapted through queries | First-class `ArrayList[T]` goal |
| Map operations | Query key-value pairs | First-class `HashMap[K, V]` goal |
| Error-aware pipelines | Not the primary model | Explicit error-aware operations are planned |

The goal is not to duplicate `go-linq`. `typed` explores a collection model that feels natural for developers moving between Go, Java, and JavaScript, while still interoperating with `iter.Seq[T]` and allowing a lazy stream when it is actually useful.

## Java / JavaScript Mapping

| Java / JavaScript | `typed` direction |
| --- | --- |
| `stream()` | `ArrayList(values).Stream()` |
| `filter` | `Filter` |
| `map` | `Map` |
| `flatMap` | `FlatMap` |
| `distinct` | `Distinct` |
| `sorted` | `Sort` / `SortBy` |
| `limit` | `Limit` |
| `skip` | `Skip` |
| `findFirst` | `First` |
| `anyMatch` | `Any` |
| `allMatch` | `All` |
| `reduce` | `Reduce` |
| `collect(toList())` | `Collect` |
| `forEach` | `ForEach` |

The project does not attempt to copy the Java or JavaScript runtime model. It borrows their collection-processing style while preserving Go's static typing, explicit errors, and straightforward control flow.

## Laziness and Execution Boundaries

A typical Stream pipeline looks like this:

```text
source -> intermediate operation -> intermediate operation -> terminal operation
```

For example:

```go
adults := stream.From(users).
    Filter(isAdult).
    Map(toProfile).
    Take(100).
    Collect()
```

Before `Collect`, `Filter`, `Map`, and `Take` only describe the pipeline. The terminal operation starts consumption, and `Take(100)` can stop the underlying source as soon as enough values have been produced.

## Error Handling

Go functions commonly return `(value, error)`, so the Stream API should not hide errors. For transformations that may fail, the project plans to support explicit error propagation, for example:

```go
profiles, err := stream.From(users).
    MapE(loadProfile).
    Collect()
```

The exact error API is still under design. The priority is to avoid silently dropping errors and to stop the pipeline promptly when an error occurs.

## Concurrency Boundary

Streams are not parallel by default. Goroutines, channels, locks, and cancellation signals have explicit concurrency semantics in Go, and automatic parallelization can introduce unpredictable overhead and lifecycle problems.

Explicit concurrent operations may be considered in the future, but ordinary `Map` will not silently become concurrent.

## Go Version

The current module uses Go 1.27:

```text
go 1.27.1
```

The target API will use Go generics and standard iterator capabilities. Go 1.23 introduced `iter.Seq`, `iter.Seq2`, and `for range` support for function iterators. Go 1.27's generic methods make fluent APIs such as `Stream[T].Map[R]` possible.

## Project Structure

```text
typed/
└── collections/
    ├── go.mod
    ├── arraylist/
    ├── hashmap/
    └── stream/
        └── mod.go
```

## Roadmap

1. Define `ArrayList[T]` and `HashMap[K, V]` semantics and naming.
2. Implement eager collection operations such as `Filter`, `Map`, `FlatMap`, and `Collect`.
3. Add `Stream()` adapters backed by `iter.Seq[T]`.
4. Add early-terminating stream operations such as `First`, `Any`, `All`, and `Take`.
5. Design error propagation, including the trade-offs around `MapE` and `FilterE`.
6. Add sorting, grouping, aggregation, and Map processing.
7. Add tests and benchmarks for eager collections, lazy streams, empty values, and early termination.

## License

This project is licensed under the [MIT License](LICENSE).
