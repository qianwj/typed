package streams

type Task func()

type Worker interface {
	Schedule(Task) Disposable
}

type Disposable interface {
	Dispose()
}

type Scheduler interface {
	Disposable
	Schedule(Task) Disposable
	CreateWorker() Worker
}
