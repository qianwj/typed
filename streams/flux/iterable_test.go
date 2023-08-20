package flux

import (
	"context"
	"fmt"
	"strconv"
	"streams"
	"testing"
)

type testIter struct {
	streams.Iterable[int]
	data []int
	idx  int
}

func (t *testIter) HasNext(context.Context) bool {
	return t.idx < len(t.data)
}

func (t *testIter) Next(context.Context) (int, error) {
	val := t.data[t.idx]
	t.idx++
	return val, nil
}

func TestIter(t *testing.T) {
	iter := &testIter{
		data: []int{1, 2, 3, 4, 5, 6},
	}
	f := FromIterable[int](iter).Map(func(a any) any {
		return strconv.Itoa(a.(int))
	}).MapT(func(a string) string {
		return a + " mapped"
	}).OnError(func(e error) {
		fmt.Printf("error: %+v\r\n", e)
	})
	f.Subscribe(func(d any) {
		fmt.Printf("data: %v\r\n", d)
	})
}
