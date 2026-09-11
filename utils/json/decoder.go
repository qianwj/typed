// Package json provides a thin wrapper over encoding/json/v2 that
// returns the project's result.Result[T] instead of the standard
// (T, error) shape.
//
// The wrapper exists for one reason: typed error handling. The
// project's other utilities (option.Optional, result.Result) keep
// the (T, error) shape but make it fluent; calling result.Success
// or result.Failure at every JSON boundary lets a pipeline read
// as one chain instead of scattering (val, err) := ...; if err != nil
// checks at every level.
//
// The package exposes two operations:
//
//   - Decode turns JSON-encoded bytes into a value of type T,
//     returning result.Result[T]. A Success carries the decoded
//     value; a Failure carries the underlying json/v2 error.
//   - Encode turns a value of type T into JSON-encoded bytes,
//     returning result.Result[[]byte]. A Success carries the
//     bytes; a Failure carries the underlying json/v2 error.
//
// Decode and Encode are inverses: a value encoded by Encode
// can be decoded back with Decode, modulo JSON-marshaling
// lossiness for channels, functions, and complex numbers.
//
// The package is not a re-export of the full encoding/json/v2
// surface. The streaming Decoder, custom Marshalers, and the
// remaining v2 features are still reachable through the
// standard import path; this package is a small surface on
// top, sized for the rest of the project.
package json

import (
	"encoding/json/v2"

	"github.com/qianwj/typed/utils/result"
)

// Decode parses the JSON-encoded data into a value of type T and
// returns it as a result.Result[T].
//
// The success path holds the decoded value:
//
//	profile, err := json.Decode[Profile](data).Unwrap()
//	if err != nil { return err }
//
// The failure path holds the underlying encoding/json/v2 error:
// a syntax error, a type mismatch, an unexpected end of input,
// or whatever the standard library surfaced. The error is not
// wrapped; callers that need to distinguish JSON errors from
// other errors can use errors.Is / errors.As against the v2
// error types.
//
// Decode is implemented in terms of result.Wrap: a single
// (value, err) pair from json.Unmarshal is wrapped into a
// Result without an explicit if-err check.
//
// # Memory
//
// Decode does not retain a reference to data after returning.
// The byte slice can be reused, pooled, or freed as soon as
// Decode returns. This is a deliberate consequence of the
// []byte (rather than string) signature: a string would have
// to be copied if the caller wanted to free the original
// buffer, but a []byte is read in place and then dropped.
//
// # Edge cases
//
// The cases below all flow through the standard encoding/json/v2
// behaviour; they are listed here so callers do not have to
// guess.
//
//   - Empty data ([]byte{}) and nil data both produce a Failure
//     with the v2 "unexpected end of JSON input" error. They do
//     not panic.
//   - Invalid JSON syntax produces a Failure with a *json.SyntaxError
//     (or the v2 equivalent) as the underlying error.
//   - A type mismatch (e.g. decoding a JSON string into an int
//     target) produces a Failure with a type-mismatch error.
//   - JSON null decoded into a non-pointer T (int, string,
//     struct, ...) produces a Success carrying the zero value
//     of T. JSON null decoded into a pointer T produces a
//     Success carrying a nil pointer. This is the v2 default
//     and matches what the v1 json package did.
//
// # Generic target
//
// T may be any Go type that the underlying json package can
// unmarshal into: basic types, structs (with field tags or
// exported fields), slices, maps with string keys, and any
// pointer to a value of those types. The wrapper does not
// impose any constraint beyond the json package's own.
func Decode[T any](data []byte) result.Result[T] {
	var r T
	err := json.Unmarshal(data, &r)
	return result.Wrap(r, err)
}
