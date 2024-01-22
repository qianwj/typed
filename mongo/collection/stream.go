package collection

import (
	"context"
	"github.com/qianwj/typed/streams"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindPublisher[V any] struct {
	cursor *mongo.Cursor
	sink   chan *streams.Item[V]
	ctx    context.Context
}

func fromCursor[V any](ctx context.Context, cursor *mongo.Cursor) streams.Publisher[V] {
	return &FindPublisher[V]{
		cursor: cursor,
		sink:   make(chan *streams.Item[V]),
		ctx:    ctx,
	}
}

func (p *FindPublisher[V]) Request(n int) {
	err := p.cursor.Err()
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if p.cursor.Next(p.ctx) {
			var val V
			err := p.cursor.Decode(val)
			p.sink <- streams.NewItem(val, err)
		} else {
			break
		}
	}
}

func (p *FindPublisher[V]) Cancel() {
	_ = p.cursor.Close(p.ctx)
	close(p.sink)
}

func (p *FindPublisher[V]) Subscribe(subscriber streams.Subscriber[V]) {
	subscriber.OnSubscribe(p)
	for item := range p.sink {
		if item.IsError() {
			subscriber.OnError(item.Err())
		} else {
			subscriber.OnNext(item.Ok())
		}
	}
	close(p.sink)
	subscriber.OnComplete()
}
