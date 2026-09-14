package objects_test

import (
	"fmt"

	"github.com/qianwj/typed/utils/objects"
)

// ExampleEquals shows the central nil-vs-empty-slice fix over
// reflect.DeepEqual: a nil slice and an empty non-nil slice compare
// equal under Equals even though reflect.DeepEqual says they differ.
func ExampleEquals() {
	var nilSlice []int
	empty := []int{}
	fmt.Println(objects.Equals(nilSlice, empty))
	fmt.Println(objects.Equals(empty, nilSlice))
	fmt.Println(objects.Equals([]int{1, 2}, []int{1, 2}))
	fmt.Println(objects.Equals([]int{1, 2}, []int{1, 3}))
	// Output:
	// true
	// true
	// true
	// false
}

// ExampleIsNil shows the typed-nil-safe nil check across pointer
// kinds that a plain `v == nil` cannot express.
func ExampleIsNil() {
	var p *int
	var iface any = p // typed nil inside an interface
	fmt.Println(objects.IsNil(iface))
	fmt.Println(objects.IsNil((*int)(nil)))
	// Output:
	// true
	// true
}
