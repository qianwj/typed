package concurrency_test

import (
	"bytes"
	"fmt"

	"github.com/qianwj/typed/concurrency"
)

func ExamplePool() {
	pool := concurrency.NewPool(func() *bytes.Buffer { return new(bytes.Buffer) })
	format := func(message string) string {
		buf := pool.Get().Get() // The creator always returns a non-nil buffer.
		buf.Reset()
		defer pool.Put(buf)
		buf.WriteString("message: ")
		buf.WriteString(message)
		return buf.String()
	}
	fmt.Println(format("hello"))
	fmt.Println(format("world"))
	// Output:
	// message: hello
	// message: world
}
