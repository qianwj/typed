package adt_test

import (
	"errors"
	"fmt"

	"github.com/qianwj/typed/adt/result"
)

// ExampleResult_Recover shows the error-aware fallback that always
// returns a T: the chain does not have to thread a (T, error) out.
func ExampleResult_Recover() {
	port := result.Wrap(lookupPort()).
		Recover(func(error) int { return 8080 })
	fmt.Println("port:", port)
	// Output: port: 8080
}

// ExampleResult_MapError shows how to enrich an error before it
// surfaces to the caller, keeping the rest of the chain typed.
func ExampleResult_MapError() {
	_, err := result.Wrap(lookupPort()).
		MapError(func(err error) error { return fmt.Errorf("config: %w", err) }).
		Unwrap()
	if err != nil {
		fmt.Println("err:", err)
	}
	// Output: err: config: not configured
}

// ExampleResult_Map demonstrates chaining a successful Result into a
// transformation while preserving error short-circuit.
func ExampleResult_Map() {
	sum := result.Success(40).
		Map(func(n int) int { return n + 2 }).
		OrElse(-1)
	fmt.Println(sum)
	// Output: 42
}

// lookupPort is a stand-in for an I/O function that fails.
func lookupPort() (int, error) { return 0, errors.New("not configured") }
