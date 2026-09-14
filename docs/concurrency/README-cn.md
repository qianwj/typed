# `typed/concurrency` — Typed

`concurrency` 是与 Go 标准库互补的并发原语包,沿用本工具集一贯的写法:具体泛型类型、不绕 `any`、不用反射。

当前版本只提供一个类型:

- `BoundedBlockingQueue[T]` —— 固定容量的 FIFO 阻塞队列,本质是 `chan T` 的一个薄泛型包装,在 channel 之上加了一层:工具集风格的命名、`Optional` 形式的非阻塞探测、带 context 的阻塞。

> 属于 **Typed** 工具集。英文原版见 [README.md](./README.md)。

## 目录

- [导入](#导入)
- [为什么需要自己实现一个队列?](#为什么需要自己实现一个队列)
- [`BoundedBlockingQueue[T]`](#boundedblockingqueuet)
  - [构造](#构造)
  - [阻塞 API](#阻塞-api)
  - [非阻塞 API](#非阻塞-api)
  - [状态查询](#状态查询)
  - [内存模型](#内存模型)
  - [与 `chan T` 的对比](#与-chan-t-的对比)
  - [示例](#示例)
- [Benchmark](#benchmark)
- [与其他包的关系](#与其他包的关系)

## 导入

```go
import (
    "github.com/qianwj/typed/concurrency"
    "github.com/qianwj/typed/utils/option"
)
```

`concurrency` 依赖 `utils/option`,因为 [`BoundedBlockingQueue.TryTake`](#非阻塞-api) 返回 `option.Optional[T]`,与工具集其他地方的“可能缺席”约定保持一致。

`concurrency` 是独立的 `go.mod` 模块,可以单独引用,不依赖 `collections` / `reactivex` / `control` / `utils` 中的任何一个。

## 为什么需要自己实现一个队列?

Go 内置的 `chan T` 本身就是一个相当不错的有界阻塞队列——前提是它有容量。Runtime 在底层用 per-P(处理器本地)队列和 lock-free 快路径实现,纯吞吐上很难被打败。事实上,`BoundedBlockingQueue` 现在就是 `chan T` 的一个薄包装:底层方法就是 `ch <- data` 和 `<-ch`,包装层每操作基本不增加任何开销。

下面这些场景下,你会想要 `BoundedBlockingQueue[T]` 套在 channel 外面:

- **非阻塞探测返回 `Optional`。** `TryTake` 返回 `option.Optional[T]`,与 `Stack.Pop` / `Queue.Pop` 同形,可以直接与工具集其他 API 链式组合,不需要再写 `(value, ok)` 风格的对偶。
- **带 context 的阻塞。** `PushCtx` / `TakeCtx` 让你直接跟超时、deadline、关停信号配合,不需要在 `Push` / `Take` 外面再自己包一层 goroutine + channel。
- **风格一致的泛型 API。** 用 `Size` 不用 `len`、用 `Capacity` 不用 `cap`、用 `Optional` 不用 `(T, bool)`——与 Typed 其他部分用同一套词汇。
- **2 的幂容量向上取整。** `Capacity()` 总是返回 2 的幂,需要做位掩码的下游用起来方便。

跟裸 `chan T` 相比,热路径上的代价基本为零。

---

## `BoundedBlockingQueue[T]`

固定容量的 FIFO 队列。容量在构造时确定,运行期不变。所有操作都对并发安全。

### 构造

```go
// NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T]
q := concurrency.NewBoundedBlockingQueue[Job](1024)
```

`capacity` 必须为正数。传 `0` 或负数会 panic —— 配错容量应当响亮地失败,而不是悄悄生成一个永远阻塞的队列。

队列的实际容量是**不小于传入值的最小 2 的幂**。`NewBoundedBlockingQueue[T](100)` 得到 `Capacity() == 128` 的队列,`NewBoundedBlockingQueue[T](1024)` 得到 `Capacity() == 1024` 的队列。这个向上取整是为了让 `Capacity()` 总是能直接当位掩码用;如果你需要精确的容量,自己传 2 的幂进来。

零值不可用,必须通过构造函数构造。

### 阻塞 API

| 方法 | 行为 |
| --- | --- |
| `Push(data T)` | 入队。**队列满时阻塞**,直到腾出空位。 |
| `Take() T` | 出队并返回。**队列空时阻塞**,直到有新元素。 |
| `PushCtx(ctx, data T) error` | 带 context 的 `Push`。阻塞到队列有空位或 `ctx` 取消;取消时返回 `ctx.Err()`,元素不入队。 |
| `TakeCtx(ctx) (T, error)` | 带 context 的 `Take`。阻塞到有元素可取或 `ctx` 取消;取消时返回 `(零值, ctx.Err())`。 |

阻塞直接走底层的 `chan T`:`Push` 就是 `ch <- data`,`Take` 就是 `<-ch`。`PushCtx` / `TakeCtx` 在同一个 `select` 上多一个 `case <-ctx.Done()`,从而零成本支持取消:ctx 取消时 `select` 走 ctx 那一支,返回 `ctx.Err()`;`PushCtx` 取消时元素不入队。

### 非阻塞 API

| 方法 | 行为 |
| --- | --- |
| `TryPush(data T) bool` | 有空位时入队,成功返回 `true`;队列满时立即返回 `false`。 |
| `TryTake() option.Optional[T]` | 队列非空时出队,成功返回 present 的 `Optional`;空队列时立即返回空 `Optional`。 |

这些方法永不等待,正好用于 `select { ... default: ... }` 风格,以及“满了就丢”或“满了就降级”这种背压策略,不需要起额外的 watcher goroutine。`TryTake` 返回 `Optional` 而非 `(T, bool)`,是工具集通用的“可能缺席”约定,与 `Stack.Pop` / `Queue.Pop` / `Deque.PopFront` / `PopBack` 一致。

### 状态查询

| 方法 | 行为 |
| --- | --- |
| `Size() int` | 当前元素数。本质是 `len(ch)`——原子读 length,不会跟紧跟着的操作串行化。窗口很小(一次原子读)但是真的;如果你需要"Size() → 紧跟着的操作看到一致视图"这种强保证,用环形缓冲区实现。与 `Stack.Size` / `Queue.Size` / `ArrayList.Size` 等命名保持一致。 |
| `Capacity() int` | 配置的容量。无锁——容量在构造后不可变。命名上对齐构造函数参数,因为 Typed 工具集里没有别的类型有固定容量。 |

### 内存模型

- **底层存储。** 启动时一次性 `make(chan T, cap)`,Go runtime 拥有 channel 内部的环形缓冲区;本类型不直接碰它。
- **2 的幂容量。** 构造函数把请求容量向上取整到 2 的幂,让 `Capacity()` 总是返回可以直接当位掩码用的值。channel 内部机制跟这个取整无关。
- **没有槽位清零。** 跟手写环形缓冲区不同,本类型不会在 `Take` / `TryTake` 时把释放的 slot 清零。带指针的元素出队后,在 channel 的底层数组里仍然存活,直到被下一次 send 覆盖。对“一次性取出很多,然后 long pause 没新 send”的工作负载,这些指针的存活时间会比带显式清零的环形缓冲区长。
- **单次操作分配。** 两种实现都为零;channel send/recv 的快路径不分配。

### `chan T` 能给你什么、不能给你什么

由于底层是 channel,跟手写 ring buffer + mutex + Cond 相比,这些 trade-off 是从 `chan T` 继承来的:

| | `chan T` | `BoundedBlockingQueue[T]` |
| --- | --- | --- |
| 容量 | 在 `make` 时设置 | 在 `NewBoundedBlockingQueue` 时设置;向上取整到 2 的幂 |
| 默认是否有界 | 否(`make(chan T)` 无缓冲) | 是——必须传 capacity |
| 阻塞发送 / 接收 | `ch <- v` / `<-ch` | `Push(v)` / `Take()` |
| 非阻塞探测 | 套 `select { default: }` | `TryPush() bool` / `TryTake() option.Optional[T]` |
| 带 context 的阻塞 | 套 `select { case <-ctx.Done(): }` | `PushCtx` / `TakeCtx` |
| `len(ch)` | 有,但与操作之间没有同步保证 | `Size()` —— 同样的语义,同样不串行化 |
| 吞吐(1P1C) | ~30 ns/op(微基准) | ~210 ns/op(微基准) |
| 吞吐(MPMC 8w) | ~22 ns/op | **~22 ns/op**(同一量级;包装层基本无开销) |

> 数字来自 `go test -bench`,Apple M5 Pro、Go 1.27,队列容量 1024(1P1C)/ 64(MPMC)。1P1C 的差距基本是微基准噪声——那个 benchmark 的消费者是 busy `TryTake` 循环,不是裸 `<-ch`,两边的开销都被它吃掉了。在生产/消费都是 uncontended `Push` / `Take` 的紧凑代码里,本类型跟裸 `chan T` 只差几 ns。它们用来确认包装层是"基本免费"的,不能替代在你的实际负载上跑一遍。

### 示例

#### 单生产者、单消费者,带背压

```go
jobs := concurrency.NewBoundedBlockingQueue[*Job](64)

go func() {
    for {
        j := jobs.Take()    // 阻塞,直到有任务
        handle(j)
    }
}()

for _, raw := range sources {
    jobs.Push(raw)          // 消费者跟不上时阻塞
}
```

`Push` 的阻塞天然就是背压:消费者慢了,生产者自动被限流,不需要额外的协议。

#### 非阻塞 offer

```go
q := concurrency.NewBoundedBlockingQueue[Event](1024)

func submit(e Event) bool {
    if !q.TryPush(e) {
        metrics.DroppedCounter.Inc()
        return false
    }
    return true
}
```

`TryPush` 可以让你直接实现“满了就丢”或“满了就降级”,不需要再起一个 watcher goroutine。

#### 查看队列状态

```go
if q.Size() > q.Capacity() * 9 / 10 {
    log.Printf("队列已用 %.0f%%", float64(q.Size()) / float64(q.Capacity()) * 100)
}
```

`Size()` 是 `len(ch)` 的一次原子 length 读,不会跟紧跟着的操作串行化。

#### 用 `TryTake` 做优雅退出

```go
stop := make(chan struct{})

go func() {
    for {
        select {
        case <-stop:
            return
        default:
        }
        if opt := q.TryTake(); !opt.IsEmpty() {
            process(opt.Get())
            continue
        }
        // 队列空:让出 CPU,让生产者跑
        runtime.Gosched()
    }
}()

// 之后停止生产者:
close(stop)
```

这种写法只通过 `TryTake` 排空队列,从不阻塞在 `Take()` 上,所以生产者结束后消费者能很快退出。

---

## Benchmark

包内自带两个 benchmark:`BoundedQueue_1P1C`、`BoundedQueue_MPMC`。运行方式:

```bash
GOWORK=off go -C concurrency test -bench=. -benchmem -run=^$ .
```

Apple M5 Pro(Go 1.27,darwin/arm64)上 1P1C 场景约 **210 ns/op**;8 worker 的 MPMC 场景约 **22 ns/op**。两者都报告 `0 B/op` 和 `0 allocs/op`。MPMC 的数字是亮点——跟裸 `chan T` 同一量级,说明包装层每操作基本无开销。1P1C 的数字主要是微基准噪声——那个 benchmark 的消费者是 busy `TryTake` 循环,不是裸 `<-ch`,两边的开销都被它吃掉了。在生产/消费都是 uncontended `Push` / `Take` 的紧凑代码里,本类型跟裸 `chan T` 只差几 ns。

## 与其他包的关系

- [`reactivex`](../../reactivex/README-cn.md) —— 带显式 demand 和可配背压的强类型异步事件流。如果要排队的其实是“派发给多个订阅者的事件”,`Observable` 通常比队列更合适。
- [`collections.Queue`](../collections/README.md#stackt--queuet--dequet) —— 同步、内存中的 `Queue[T]`。在没有并发、又想要 `Optional[T]` 风格取值时使用;它不加锁、没有 `TryPush`、也没有背压。
- 标准库的 [`chan T`](https://go.dev/ref/spec#Channel_types) —— 本类型就是它的一个薄包装。
