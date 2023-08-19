package flux

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestJust(t *testing.T) {
	f := Just(1, 2, errors.New("an error"), 3, 4, 5).Map(func(a any) any {
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
