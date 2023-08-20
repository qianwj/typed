package streams

import "context"

type Iterable[T any] interface {
	HasNext(ctx context.Context) bool
	Next(ctx context.Context) (T, error)
}
