package objects

import "testing"

func TestIsNil(t *testing.T) {
	var (
		nilPointer *int
		nilMap     map[string]int
		nilSlice   []int
		nilChan    chan int
		nilFunc    func()
		nilIface   any
	)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{name: "untyped nil", value: nil, want: true},
		{name: "nil pointer", value: nilPointer, want: true},
		{name: "nil map", value: nilMap, want: true},
		{name: "nil slice", value: nilSlice, want: true},
		{name: "nil channel", value: nilChan, want: true},
		{name: "nil function", value: nilFunc, want: true},
		{name: "nil interface", value: nilIface, want: true},
		{name: "pointer", value: new(int), want: false},
		{name: "map", value: map[string]int{}, want: false},
		{name: "slice", value: []int{}, want: false},
		{name: "channel", value: make(chan int), want: false},
		{name: "function", value: func() {}, want: false},
		{name: "integer", value: 0, want: false},
		{name: "string", value: "", want: false},
		{name: "struct", value: struct{}{}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsNil(test.value); got != test.want {
				t.Fatalf("IsNil(%s): got %v, want %v", test.name, got, test.want)
			}
		})
	}
}
