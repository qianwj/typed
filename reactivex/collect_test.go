package reactivex

import (
	"context"
	"reflect"
	"testing"
)

func TestToSliceAndCollect(t *testing.T) {
	values, err := FromSlice([]int{1, 2, 3}).ToSlice(context.Background())
	if err != nil || !reflect.DeepEqual(values, []int{1, 2, 3}) {
		t.Fatalf("values=%v err=%v", values, err)
	}
	total, err := Collect(context.Background(), FromSlice([]int{1, 2, 3}), 0, func(sum, value int) int { return sum + value })
	if err != nil || total != 6 {
		t.Fatalf("total=%d err=%v", total, err)
	}
}
