package control

import (
	"slices"
	"testing"
)

func TestIfEvaluatesBothArguments(t *testing.T) {
	for _, tc := range []struct {
		name      string
		condition bool
		want      int
	}{
		{"true", true, 0},
		{"false", false, 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			trueCalls, falseCalls := 0, 0
			got := If(tc.condition,
				func() int { trueCalls++; return 0 }(),
				func() int { falseCalls++; return 42 }(),
			)
			if got != tc.want || trueCalls != 1 || falseCalls != 1 {
				t.Fatalf("got=%d trueCalls=%d falseCalls=%d, want %d, 1, 1", got, trueCalls, falseCalls, tc.want)
			}
		})
	}
}

func TestIfSupportsNonComparableValues(t *testing.T) {
	if got := If(true, []int{1, 2}, []int{3}); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v, want [1 2]", got)
	}
}

func TestIfGetCallsOnlySelectedBranch(t *testing.T) {
	for _, tc := range []struct {
		name      string
		condition bool
		want      []int
		trueCalls int
	}{
		{"true", true, []int{1, 2}, 1},
		{"false", false, nil, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			trueCalls, falseCalls := 0, 0
			got := IfGet(tc.condition,
				func() []int { trueCalls++; return []int{1, 2} },
				func() []int { falseCalls++; return nil },
			)
			if !slices.Equal(got, tc.want) || trueCalls != tc.trueCalls || falseCalls != 1-tc.trueCalls {
				t.Fatalf("got=%v trueCalls=%d falseCalls=%d", got, trueCalls, falseCalls)
			}
		})
	}
}

func TestIfGetAllowsUnselectedNilCallback(t *testing.T) {
	value := func() string { return "selected" }
	if got := IfGet(true, value, nil); got != "selected" {
		t.Fatalf("true branch: got %q", got)
	}
	if got := IfGet(false, nil, value); got != "selected" {
		t.Fatalf("false branch: got %q", got)
	}
}

func TestIfGetSelectedNilCallbackPanics(t *testing.T) {
	for _, name := range []string{"true", "false"} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected selected nil callback to panic")
				}
			}()
			IfGet[int](name == "true", nil, nil)
		})
	}
}
