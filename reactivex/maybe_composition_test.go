package reactivex

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestMaybeCompositionStatesAndLaziness(t *testing.T) {
	t.Parallel()
	number := 42
	leftErr, rightErr := errors.New("left"), errors.New("right")
	type state struct {
		name    string
		value   *int
		present bool
		err     error
	}
	states := []state{
		{name: "value", value: &number, present: true},
		{name: "nil", present: true},
		{name: "empty", value: &number},
		{name: "error", value: &number, present: true, err: leftErr},
	}
	for _, operator := range []string{"Map", "FlatMap", "Zip", "AndThen"} {
		for _, left := range states {
			for _, right := range states {
				t.Run(operator+"/"+left.name+"/"+right.name, func(t *testing.T) {
					if right.err != nil {
						right.err = rightErr
					}
					var leftCalls, rightCalls, transformCalls atomic.Int32
					source := NewMaybe(func() (*int, bool, error) {
						leftCalls.Add(1)
						return left.value, left.present, left.err
					})
					next := NewMaybe(func() (*int, bool, error) {
						rightCalls.Add(1)
						return right.value, right.present, right.err
					})
					var composed *Maybe[*int]
					switch operator {
					case "Map":
						composed = source.Map(func(v *int) *int {
							transformCalls.Add(1)
							return v
						})
					case "FlatMap":
						composed = source.FlatMap(func(*int) *Maybe[*int] {
							transformCalls.Add(1)
							return next
						})
					case "Zip":
						composed = source.Zip(next, func(_, v *int) *int {
							transformCalls.Add(1)
							return v
						})
					case "AndThen":
						composed = source.AndThen(next)
					}
					if leftCalls.Load() != 0 || rightCalls.Load() != 0 || transformCalls.Load() != 0 {
						t.Fatal("composition must not execute sources or transforms before Await")
					}
					want := left
					var wantRight, wantTransform int32
					if left.err == nil && left.present {
						if operator != "Map" {
							want, wantRight = right, 1
						}
						if operator == "Map" || operator == "FlatMap" || (operator == "Zip" && right.err == nil && right.present) {
							wantTransform = 1
						}
					}
					for range 2 {
						value, err := composed.Await().Unwrap()
						if !errors.Is(err, want.err) || value.IsPresent() != (want.err == nil && want.present) {
							t.Fatalf("Await = (%v, %v), want presence %v, error %v", value, err, want.err == nil && want.present, want.err)
						}
						if value.IsPresent() && value.Get() != want.value {
							t.Fatalf("value = %v, want %v", value.Get(), want.value)
						}
					}
					if leftCalls.Load() != 1 || rightCalls.Load() != wantRight || transformCalls.Load() != wantTransform {
						t.Fatalf("calls (left, right, transform) = (%d, %d, %d), want (1, %d, %d)", leftCalls.Load(), rightCalls.Load(), transformCalls.Load(), wantRight, wantTransform)
					}
				})
			}
		}
	}
}
