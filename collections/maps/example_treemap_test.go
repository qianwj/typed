package maps_test

import (
	"cmp"
	"fmt"

	"github.com/qianwj/typed/collections/maps"
)

func ExampleTreeMap() {
	m := maps.NewTreeMap[int, string](cmp.Compare[int])
	m.Put(30, "thirty")
	m.Put(10, "ten")
	m.Put(20, "twenty")
	fmt.Println(m.Keys().Collect())
	fmt.Println(m.Floor(25).Get())
	fmt.Println(m.Higher(20).Get())
	fmt.Println(m.Range(10, 30).Collect())
	fmt.Println(m.Stream().Filter(func(k int, _ string) bool { return k == 20 }).
		MapValues(func(_ int, v string) int { return len(v) }).Collect())
	// Output:
	// [10 20 30]
	// {20 twenty}
	// {30 thirty}
	// [{10 ten} {20 twenty}]
	// map[20:6]
}
