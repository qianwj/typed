# `typed/concurrency` — Typed

`concurrency` 是与 Go 标准库互补的并发原语包,沿用本工具集一贯的写法:具体泛型类型、不绕 `any`、不用反射。

当前版本只提供一个类型:

- `BoundedBlockingQueue[T]` —— 固定容量的 FIFO 队列,带阻塞与非阻塞两套 API,内部用单个环形缓冲区 + 单个 `sync.Mutex` + 单个 `*sync.Cond` 实现。

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

Go 内置的 `chan T` 本身就是一个相当不错的有界阻塞队列——前提是它有容量。Runtime 在底层用 per-P(处理器本地)队列和 lock-free 快路径实现,纯吞吐上很难被打败。

下面这些场景下,`chan T` 的表达力不够,你会想要 `BoundedBlockingQueue[T]`:

- **非阻塞探测。** `TryPush` / `TryTake` 让你“先看一眼、再决定”而不必为 `select` 起一个超时的 goroutine。`TryTake` 返回 `option.Optional[T]`,与 `Stack.Pop` / `Queue.Pop` 同形,可以与工具集其他 API 链式组合。
- **同步的状态观测。** `Size()` / `Capacity()` 在与入队/出队同一把锁下取值,你可以基于背压做条件分支而不用单独维护一个计数器。
- **风格一致的泛型 API。** 整个仓库已经在用 Typed 工具集时,这个队列的签名就跟它们长在一起:用 `Size` 不用 `len`、用 `Optional` 不用 `(T, bool)`。

代价是吞吐。在 1P1C 的微基准上,原生 `chan` 大约快 **5–7 倍**。数组环形缓冲区把差距压到了不那么大——两种实现都是**每次操作零分配**,所以 GC 压力不是分水岭,同步原语的开销才是。

`BoundedBlockingQueue` 的定位因此是:**“`chan T` 表达力不够,又能接受吞吐税时使用”**。

---

## `BoundedBlockingQueue[T]`

固定容量的 FIFO 队列。容量在构造时确定,运行期不变。所有操作都对并发安全。

### 构造

```go
// NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T]
q := concurrency.NewBoundedBlockingQueue[Job](1024)
```

`capacity` 必须为正数。传 `0` 或负数会 panic —— 配错容量应当响亮地失败,而不是悄悄生成一个永远阻塞的队列。

队列的实际容量是**不小于传入值的最小 2 的幂**。`NewBoundedBlockingQueue[T](100)` 得到 `Capacity() == 128` 的队列,`NewBoundedBlockingQueue[T](1024)` 得到 `Capacity() == 1024` 的队列。这个向上取整的策略让 `Push` / `Take` 在推进 head / tail 时用 `& mask` 一次位运算就能替代 `% cap`,这是这条策略的动机。如果你需要精确的容量,自己传 2 的幂进来。

零值不可用,必须通过构造函数构造。

### 阻塞 API

| 方法 | 行为 |
| --- | --- |
| `Push(data T)` | 入队。**队列满时阻塞**,直到腾出空位。 |
| `Take() T` | 出队并返回。**队列空时阻塞**,直到有新元素。 |

阻塞通过一个与 mutex 共享的 `*sync.Cond` 实现。每次状态变化(`Push` / `Take` / `TryPush` / `TryTake` / `DrainTo`)都调用 `cond.Broadcast()`——这是相对保守的选择:总能正确唤醒需要的等待者(`Take` 或 `DrainTo` 之后唤醒 pushers,`Push` 之后唤醒 takers),不需要拆成两个 `Cond` 也不需要怕 thundering-herd 的 `Signal()`。本来也只有一把锁,这个复杂度刚刚好。

### 非阻塞 API

| 方法 | 行为 |
| --- | --- |
| `TryPush(data T) bool` | 有空位时入队,成功返回 `true`;队列满时立即返回 `false`。 |
| `TryTake() option.Optional[T]` | 队列非空时出队,成功返回 present 的 `Optional`;空队列时立即返回空 `Optional`。 |
| `DrainTo(dst []T) int` | 出队最多 `len(dst)` 个元素,按 FIFO 顺序写入 `dst`,返回实际写入的个数。`dst` 不会被扩容;如果队列元素少于 `len(dst)`,`dst` 中剩余的槽位保持不动。 |

这些方法永不等待,正好用于 `select { ... default: ... }` 风格,以及“满了就丢”或“满了就降级”这种背压策略,不需要起额外的 watcher goroutine。`TryTake` 返回 `Optional` 而非 `(T, bool)`,是工具集通用的“可能缺席”约定,与 `Stack.Pop` / `Queue.Pop` / `Deque.PopFront` / `PopBack` 一致。

`DrainTo` 的存在是为了把一批出队动作放在同一个临界区里完成——逐个 `TryTake` 需要 N 次取锁,且这 N 个元素之间不构成一个一致的快照;`DrainTo` 一次性完成。

### 状态查询

| 方法 | 行为 |
| --- | --- |
| `Size() int` | 当前元素数。取用与入队/出队相同的锁,所以返回值与紧随其后的 `Push` / `Take` 之间是一致的快照。与 `Stack.Size` / `Queue.Size` / `ArrayList.Size` 等保持一致。 |
| `Capacity() int` | 配置的容量。无锁——容量在构造后不可变。命名上对齐构造函数参数,因为 Typed 工具集里没有别的类型有固定容量。 |

### 内存模型

- **底层存储。** 启动时一次性 `make([]T, capacity)`,head 与 tail 索引在它上面环绕。热路径上无任何每元素分配。
- **2 的幂容量。** 构造函数把请求容量向上取整到 2 的幂,并把 `cap - 1` 预存为 `mask`。`head` / `tail` 推进时用 `& mask`——一次位运算——替代 `% cap`。空 vs 满由 `count` 字段区分,所以 head 等于 tail 不会造成歧义。
- **槽位清零。** `Take` / `TryTake` / `DrainTo` 在推进 `head` 之前,会用 `T` 的零值覆盖刚刚释放的槽位。这只是给 GC 的提示:如果 `T` 含有指针,出队的值在队列还持有底层数组时也能被回收。该槽位在下次 `Push` 覆写前都不会再被读。
- **内部状态不外泄指针。** slice 由结构体按值持有,所有调用者都走同一把 mutex。

### 与 `chan T` 的对比

| | `chan T` | `BoundedBlockingQueue[T]` |
| --- | --- | --- |
| 容量 | 在 `make` 时设置 | 在 `NewBoundedBlockingQueue` 时设置;向上取整到 2 的幂 |
| 默认是否有界 | 否(`make(chan T)` 无缓冲) | 是——必须传 capacity |
| 阻塞发送 / 接收 | `ch <- v` / `<-ch` | `Push(v)` / `Take()` |
| 非阻塞探测 | 无——要套 `select { default: }` | `TryPush() bool` / `TryTake() option.Optional[T]` |
| 批量出队 | 无——循环里逐个取 | `DrainTo(dst []T) int` |
| `len(ch)` | 有,但与操作之间没有同步保证 | `Size()`,与操作在同一把锁下 |
| 单次操作分配 | 0 | 0 |
| 吞吐(1P1C) | ~30 ns/op | ~200 ns/op |
| 吞吐(MPMC 8w) | ~22 ns/op | ~118 ns/op |
| 批量出队(64) | n/a | ~310 ns/op(零分配) |

> 数字来自 `go test -bench`,Apple M4 Pro、Go 1.27.1,队列容量 1024(1P1C)/ 64(MPMC)。它们只是用来确认这个队列与手写的 `sync.Cond` 实现在同一量级,不能替代在你的实际负载上跑一遍。

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

`Size()` 在与 `Push` / `Take` 同一把锁下取值,所以读到的值与紧随其后的操作是一致的。

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

#### 批量 drain

如果你想把一批数据放在同一个临界区里处理——比如 worker 在“刷盘 tick”时一次性 flush——`DrainTo` 比循环 `TryTake` 更便宜、且更一致:

```go
buf := make([]Event, 64) // 按你期望的 batch 大小确定
for {
    n := q.DrainTo(buf)
    if n == 0 {
        time.Sleep(idleInterval) // 或者在 Take 上阻塞,如果不关心关停的话
        continue
    }
    flushBatch(buf[:n])
}
```

两点注意:

- `DrainTo` 不会扩容 `dst` —— 如果队列中元素超过 `len(dst)`,剩余的会留在队列里。把 `len(dst)` 设成你实际想处理的 batch 大小即可。
- `DrainTo` 会调用 `cond.Broadcast()`,所以任何阻塞在满队列上的 `Push` 都会在 drain 完成后立刻被唤醒。

---

## Benchmark

包内自带两组对照 benchmark:`BoundedBlockingQueue_*` 与等价的 `chan int` 基线。运行方式:

```bash
GOWORK=off go -C concurrency test -bench=. -benchmem -run=^$ .
```

Apple M4 Pro(Go 1.27.1,darwin/arm64)上 1P1C 场景下,`BoundedBlockingQueue` 落在 **195 ns/op** 左右,`chan int` 约 **28 ns/op**;8 worker 的 MPMC 场景下,前者约 **118 ns/op**,后者约 **22 ns/op**。两种实现都报告 `0 B/op` 和 `0 allocs/op`。

这些 benchmark 的意义在于确认基于数组的设计不会在 GC 压力上退化。它们不能替代你在自己负载上的测量。

## 未来工作 取消支持

`BoundedBlockingQueue` 目前没有 `context.Context` 感知的 `Push` / `Take` 变体。本节把设计思路写下来,避免下一个人需要取消时还要重新调研一遍。

### 拟定的签名

```go
// PushCtx 阻塞,直到有空位或 ctx 被取消。
// 入队成功返回 nil;被取消时不入队,返回 ctx.Err()。
func (q *BoundedBlockingQueue[T]) PushCtx(ctx context.Context, data T) error

// TakeCtx 阻塞,直到队列非空或 ctx 被取消。
// 取到元素时返回 (v, nil);被取消时返回 (零值, ctx.Err())。
func (q *BoundedBlockingQueue[T]) TakeCtx(ctx context.Context) (T, error)
```

这套签名对齐 Go 既有的约定(`http.Request.WithContext`、`sql.QueryContext`,`sync.WaitGroup` 等也都通过 channel 跟 ctx 配合),让调用方在超时、deadline、关停信号这些场景下不用自己在 `Push` / `Take` 外面再包一层 goroutine + channel 的脚手架。

### 卡点:`sync.Cond` 不支持 `context.Context`

`Push` 和 `Take` 用 `sync.Cond.Wait()` 把 goroutine 挂起,等 `Broadcast()`。`Cond.Wait` 是“挂在 mutex 保护的条件变量上”最干净的写法,但它没有取消钩子——走出 `cond.Wait` 的唯一方式就是 `Broadcast` / `Signal`,没法 `select` 它。所以要加 `PushCtx` / `TakeCtx`,必须二选一:

**方案 A —— 每次调用起一个 watcher goroutine。**

每次 `PushCtx` 都起一个 goroutine 干这件事: `<-ctx.Done(); q.cond.Broadcast()` 然后退出。`Broadcast` 本身很便宜,但每个阻塞调用都额外多一个 goroutine 加一个 channel,会体现在 `pprof` 里,也会在高负载下变成分配器压力。偶尔用用还行,作为默认路径就太浪费了。

**方案 B —— 把 `sync.Cond` 换成基于 channel 的信号。**

队列里多两个 `cap=1` 的 `chan struct{}`:`notEmpty` 与 `notFull`。`Push` / `PushCtx` 往 `notEmpty` 发;`Take` / `TakeCtx` 往 `notFull` 发。阻塞就变成 `select { <-signal; <-ctx.Done() }`。Go 自己的 `sync/semaphore` 就是这么做的,对 ctx 感知 API 来说这是正确形态。

代价:基本 `Push` / `Take` 每次操作要发一次 channel,而之前是 `Cond.Broadcast`(`cap=1` 的 buffered channel 在槽位空时发送几乎免费,稳态下就是一次缓存行写入)。更大的代价是改写——`Push` / `Take` / `TryPush` / `TryTake` / `DrainTo` 全部和等待队列有交互,「先取锁、再发信号、再放锁」的顺序要仔细推敲,不然会漏唤醒或重现 thundering-herd。

### 为什么暂缓

目前这个包还没有 ctx 变体的真实调用方。在没有真实需求的情况下:

- 走方案 A,会让以后读这份代码的人都要问一句“为什么每次调用多一个 goroutine”;
- 走方案 B 是正确答案,但改动够大,如果在没有真实用例的情况下做,容易在 broadcast 顺序上踩到只有 contention 高的测试才能复现的坑。

计划是在真正需要时,把同步核心整体改成方案 B,补一个“多个 pushers、消费者中途取消、所有 pushers 都观察到取消、且不漏 broadcast”的测试,然后再把 `PushCtx` / `TakeCtx` 公开。在此之前,需要取消的调用方可以自己在 goroutine 里包一层 `Push` / `Take`,用 `select` 配 `ctx.Done()`。

## 与其他包的关系

- [`reactivex`](../../reactivex/README-cn.md) —— 带显式 demand 和可配背压的强类型异步事件流。如果要排队的其实是“派发给多个订阅者的事件”,`Observable` 通常比队列更合适。
- [`collections.Queue`](../collections/README.md#stackt--queuet--dequet) —— 同步、内存中的 `Queue[T]`。在没有并发、又想要 `Optional[T]` 风格取值时使用;它不加锁、没有 `TryPush`、也没有背压。
- 标准库的 [`sync.Cond`](https://pkg.go.dev/sync#Cond) 与 [`chan T`](https://go.dev/ref/spec#Channel_types) —— 本类型的底层基石。
