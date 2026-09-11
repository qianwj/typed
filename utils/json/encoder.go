package json

import (
	"encoding/json/v2"

	"github.com/qianwj/typed/utils/result"
)

// Encode serialises a value of type T into a JSON byte slice and
// returns it as a result.Result[[]byte].
//
// The success path holds the JSON-encoded bytes:
//
//	data, err := json.Encode(profile).Unwrap()
//	if err != nil { return err }
//
// The failure path holds the underlying encoding/json/v2 error:
// an unencodable value (a function, a channel, ...), a custom
// Marshaler that returned a non-nil error, or anything the
// standard library surfaced. The error is not wrapped; callers
// that need to distinguish JSON errors from other errors can
// use errors.Is / errors.As against the v2 error types.
//
// Encode is implemented in terms of result.Wrap: a single
// ([]byte, error) pair from json.Marshal is wrapped into a
// Result without an explicit if-err check.
//
// # Options
//
// The opts parameter is a variadic of encoding/json/v2 Options
// (jsonopts.Options). Each option is a property setter; later
// options override earlier ones. Common cases:
//
//	// Always produce the same byte sequence for the same
//	// input (sorts map keys, normalises number formatting).
//	data := json.Encode(profile, json.Deterministic(true)).OrElse(nil)
//
//	// Skip fields whose value is the type's zero value.
//	data := json.Encode(profile, json.OmitZeroStructFields(true)).OrElse(nil)
//
// v2 options are different from v1 json.Marshal's tags like
// MarshalJSON / UnmarshalJSON. v2's options are runtime
// configuration; tags and Marshaler methods are still the way
// to influence what is produced.
//
// # Round-trip with Decode
//
// Encode and Decode are inverses. A value encoded by Encode
// and then decoded back into the same target type with Decode
// produces an equal value (by Go's == comparison on the
// decoded fields, modulo the JSON-marshaling lossiness for
// channels, functions, and complex numbers). The round-trip
// property is the contract the rest of the project relies
// on; the tests in encoder_test.go cover it explicitly.
//
// # Memory
//
// Encode allocates a fresh []byte for the encoded value. The
// caller owns the returned slice and may reuse, pool, or
// free it as soon as the Result is consumed. The input value
// t is read by json.Marshal but not retained; it can be
// mutated as soon as Encode returns.
//
// # HTML escaping
//
// Unlike encoding/json v1, encoding/json/v2 does not escape
// HTML by default. A value containing <, >, or & is encoded
// verbatim. If HTML escaping is required, callers should run
// the bytes through a separate escaping step or use the
// v1-style options explicitly.
//
// # Type constraint
//
// T may be any Go type that the underlying json package can
// marshal: basic types, structs (with field tags or exported
// fields), slices, maps with string keys, and any pointer
// to a value of those types. A pointer to a nil T marshals to
// "null". An unencodable value (function, channel, complex
// number) returns a Failure with a *json.UnsupportedTypeError
// or the v2 equivalent.
func Encode[T any](t T, opts ...json.Options) result.Result[[]byte] {
	return result.Wrap(json.Marshal(t, opts...))
}
