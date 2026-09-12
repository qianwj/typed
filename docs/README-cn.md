# `typed` 文档

按包组织的 API 参考和示例。每篇都按"导入 → 设计动机 → API → 例子 → 与其他包的关系"五段排版，可以独立读。

> Looking for the English version? See [README.md](./README.md).

## 集合与同步流

- [`collections`](./collections/README-cn.md) — `Stack` / `Queue` / `Deque` / `ArrayList` / `LinkedList` / `HashMap` / `HashSet` / `Stream` / `Range`，以及 `collections.Range`。
- [`reactivex`](./reactivex/README-cn.md) — 异步事件流：`Observable` / `Publisher` / `Subscriber` / `Subject`、背压、源与算子。

## 控制流

- [`control`](./control/README-cn.md) — `Repeat` / `RepeatE` 与 `control/match` 模式匹配（`Pattern[T]` + 链式 `Case` / `Type`）。

## 抽象

- [`utils/option`](./option/README-cn.md) — `Optional[T]`，可能缺席的值。
- [`utils/result`](./result/README-cn.md) — `Result[T]`，成功 / 失败以及与 `(T, error)` 互转的桥。
- [`utils/objects`](./utils/objects/README-cn.md) — `IsNil` / `Equals`（含 `Equal(T) bool` 分派与 nil/空归一化）。
- [`utils/json`](./utils/json/README-cn.md) — `Encode` / `Decode`，基于 `encoding/json/v2` 的 `Result` 风格编解码。

## 阅读建议

- 第一次接触：从 [`collections`](./collections/README-cn.md) 看 `ArrayList` / `LinkedList` / `Stream` 的不可变 transform 风格，再到 [`utils/option`](./option/README-cn.md) / [`utils/result`](./result/README-cn.md) 看缺席 / 失败的表达。
- 需要异步或多订阅：跳到 [`reactivex`](./reactivex/README-cn.md)。
- 需要按值或类型做 dispatch：跳到 [`control`](./control/README-cn.md)。
