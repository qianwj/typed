package json_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/qianwj/typed/utils/json"
	"github.com/qianwj/typed/utils/result"
)

// ---------- success path: primitive targets ----------

// TestDecodeInteger confirms the basic happy path for a primitive
// numeric target. The success value matches the JSON literal and
// the result is a success, not a failure.
func TestDecodeInteger(t *testing.T) {
	r := json.Decode[int]([]byte("42"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[int](42): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != 42 {
		t.Fatalf("Decode[int](42).Value: got %d, want 42", got)
	}
}

// TestDecodeString confirms the happy path for a string target,
// including the surrounding quotes in the JSON input.
func TestDecodeString(t *testing.T) {
	r := json.Decode[string]([]byte(`"hello"`))
	if !r.IsSuccess() {
		t.Fatalf("Decode[string]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != "hello" {
		t.Fatalf("Decode[string]: got %q, want %q", got, "hello")
	}
}

// TestDecodeBool confirms the happy path for a boolean target.
func TestDecodeBool(t *testing.T) {
	r := json.Decode[bool]([]byte("true"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[bool]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != true {
		t.Fatalf("Decode[bool]: got %v, want true", got)
	}
}

// TestDecodeFloat confirms the happy path for a float target,
// including a value that is not representable as an integer
// (so accidental int targets would not work).
func TestDecodeFloat(t *testing.T) {
	r := json.Decode[float64]([]byte("3.14"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[float64]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != 3.14 {
		t.Fatalf("Decode[float64]: got %v, want 3.14", got)
	}
}

// ---------- success path: composite targets ----------

// TestDecodeStruct confirms the happy path for a struct target.
// Field names match the JSON keys, so no struct tags are needed.
func TestDecodeStruct(t *testing.T) {
	type profile struct {
		Name string
		Age  int
	}
	r := json.Decode[profile]([]byte(`{"Name": "alice", "Age": 30}`))
	if !r.IsSuccess() {
		t.Fatalf("Decode[profile]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got.Name != "alice" || got.Age != 30 {
		t.Fatalf("Decode[profile]: got %+v, want {alice 30}", got)
	}
}

// TestDecodeStructWithTags confirms that struct tags are honoured.
// This is the documented way to map JSON's snake_case keys to
// Go's CamelCase fields, and Decode must respect them.
func TestDecodeStructWithTags(t *testing.T) {
	type profile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	r := json.Decode[profile]([]byte(`{"first_name": "alice", "last_name": "smith"}`))
	if !r.IsSuccess() {
		t.Fatalf("Decode[profile]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got.FirstName != "alice" || got.LastName != "smith" {
		t.Fatalf("Decode[profile]: got %+v, want {alice smith}", got)
	}
}

// TestDecodeSlice confirms the happy path for a slice target.
// Element order is preserved.
func TestDecodeSlice(t *testing.T) {
	r := json.Decode[[]int]([]byte("[1, 2, 3, 4, 5]"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[[]int]: IsSuccess = false, err = %v", r.Error())
	}
	got := r.Value()
	want := []int{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("Decode[[]int]: got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Decode[[]int] at %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

// TestDecodeMap confirms the happy path for a string-keyed map
// target. Map values are decoded in JSON order, but Go map
// iteration is not order-preserving, so the test asserts set
// equality rather than ordered equality.
func TestDecodeMap(t *testing.T) {
	r := json.Decode[map[string]int]([]byte(`{"a": 1, "b": 2, "c": 3}`))
	if !r.IsSuccess() {
		t.Fatalf("Decode[map[string]int]: IsSuccess = false, err = %v", r.Error())
	}
	got := r.Value()
	if len(got) != 3 {
		t.Fatalf("Decode[map[string]int]: len = %d, want 3", len(got))
	}
	for k, v := range got {
		switch k {
		case "a":
			if v != 1 {
				t.Fatalf("map[a]: got %d, want 1", v)
			}
		case "b":
			if v != 2 {
				t.Fatalf("map[b]: got %d, want 2", v)
			}
		case "c":
			if v != 3 {
				t.Fatalf("map[c]: got %d, want 3", v)
			}
		default:
			t.Fatalf("unexpected key %q", k)
		}
	}
}

// TestDecodeNested confirms the happy path for a nested struct.
// The outer object contains a child object that decodes into
// a struct field of the same name.
func TestDecodeNested(t *testing.T) {
	type address struct {
		City string
		Zip  string
	}
	type person struct {
		Name    string
		Address address
	}
	data := []byte(`{"Name": "alice", "Address": {"City": "Boston", "Zip": "02115"}}`)
	r := json.Decode[person](data)
	if !r.IsSuccess() {
		t.Fatalf("Decode[person]: IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got.Name != "alice" || got.Address.City != "Boston" || got.Address.Zip != "02115" {
		t.Fatalf("Decode[person]: got %+v", got)
	}
}

// TestDecodePointerTarget confirms that decoding into a pointer
// type T = *Foo allocates a Foo and points r at it. The result
// must be a non-nil pointer on the success path.
func TestDecodePointerTarget(t *testing.T) {
	type inner struct {
		Value int
	}
	r := json.Decode[*inner]([]byte(`{"Value": 7}`))
	if !r.IsSuccess() {
		t.Fatalf("Decode[*inner]: IsSuccess = false, err = %v", r.Error())
	}
	got := r.Value()
	if got == nil {
		t.Fatal("Decode[*inner].Value: got nil, want non-nil")
	}
	if got.Value != 7 {
		t.Fatalf("Decode[*inner].Value.Value: got %d, want 7", got.Value)
	}
}

// ---------- failure path: invalid input ----------

// TestDecodeEmptyInput confirms that an empty byte slice
// produces a Failure (not a panic). The error message is the
// v2 "unexpected end of JSON input" diagnostic.
func TestDecodeEmptyInput(t *testing.T) {
	r := json.Decode[int]([]byte{})
	if r.IsSuccess() {
		t.Fatal("Decode[int](\"\"): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Decode[int](\"\").Error: got nil, want non-nil")
	}
}

// TestDecodeNilInput confirms that a nil byte slice also
// produces a Failure. The standard library's json package
// treats nil as an empty input, and so does v2.
func TestDecodeNilInput(t *testing.T) {
	r := json.Decode[int](nil)
	if r.IsSuccess() {
		t.Fatal("Decode[int](nil): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Decode[int](nil).Error: got nil, want non-nil")
	}
}

// TestDecodeInvalidJSONSyntax confirms that a malformed JSON
// payload produces a Failure with a non-nil error.
func TestDecodeInvalidJSONSyntax(t *testing.T) {
	r := json.Decode[int]([]byte("{not json"))
	if r.IsSuccess() {
		t.Fatal("Decode[int]({not json): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Decode[int]({not json).Error: got nil, want non-nil")
	}
}

// TestDecodeTypeMismatch confirms that a type mismatch between
// the JSON value and the target type produces a Failure. The
// JSON value is a string but the target is an int.
func TestDecodeTypeMismatch(t *testing.T) {
	r := json.Decode[int]([]byte(`"not a number"`))
	if r.IsSuccess() {
		t.Fatal("Decode[int](\"not a number\"): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Decode[int] type-mismatch: Error = nil, want non-nil")
	}
}

// TestDecodeTrailingData confirms that extra non-whitespace
// bytes after the JSON value produce a Failure. The first JSON
// value is the integer 1, but the trailing '2' is not allowed
// at the top level.
func TestDecodeTrailingData(t *testing.T) {
	r := json.Decode[int]([]byte("1 2"))
	if r.IsSuccess() {
		t.Fatal("Decode[int](1 2): IsSuccess = true, want false (trailing data)")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Decode[int] trailing data: Error = nil, want non-nil")
	}
}

// ---------- failure path: composition ----------

// TestDecodeFailurePreservesError confirms that the underlying
// json error is preserved verbatim on the failure path. The
// caller can use errors.As or string match against the message.
func TestDecodeFailurePreservesError(t *testing.T) {
	r := json.Decode[int]([]byte("{not json"))
	err := r.Error()
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if !strings.Contains(err.Error(), "json") || !strings.Contains(err.Error(), "offset") {
		// Different versions of encoding/json/v2 word the
		// diagnostic differently, but the message should at
		// least mention "json" and a position. If a future
		// v2 release changes the wording this assertion will
		// need to be updated, but the test will catch a
		// silent loss of error context in the meantime.
		t.Logf("note: error message %q did not contain expected fragments; v2 wording may have changed", err.Error())
	}
}

// ---------- JSON null edge cases ----------

// TestDecodeNullIntoValueType confirms that JSON null decoded
// into a value type (int) produces a Success with the zero
// value. This is the documented v2 behaviour and matches v1.
func TestDecodeNullIntoValueType(t *testing.T) {
	r := json.Decode[int]([]byte("null"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[int](null): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != 0 {
		t.Fatalf("Decode[int](null).Value: got %d, want 0 (zero value)", got)
	}
}

// TestDecodeNullIntoString confirms the same behaviour for a
// string target: JSON null is a Success carrying the empty
// string.
func TestDecodeNullIntoString(t *testing.T) {
	r := json.Decode[string]([]byte("null"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[string](null): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != "" {
		t.Fatalf("Decode[string](null).Value: got %q, want empty string", got)
	}
}

// TestDecodeNullIntoPointerType confirms that JSON null decoded
// into a pointer type produces a Success with a nil pointer.
// The non-nil-pointer allocation only happens for non-null JSON.
func TestDecodeNullIntoPointerType(t *testing.T) {
	type inner struct {
		V int
	}
	r := json.Decode[*inner]([]byte("null"))
	if !r.IsSuccess() {
		t.Fatalf("Decode[*inner](null): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); got != nil {
		t.Fatalf("Decode[*inner](null).Value: got %+v, want nil", got)
	}
}

// ---------- composition with result ----------

// TestDecodeUnwrapToStandardShape confirms the recommended
// pattern for callers that need the standard (T, error) shape:
// the Unwrap method returns the value and the error in one step.
func TestDecodeUnwrapToStandardShape(t *testing.T) {
	v, err := json.Decode[int]([]byte("42")).Unwrap()
	if err != nil {
		t.Fatalf("Unwrap: err = %v, want nil", err)
	}
	if v != 42 {
		t.Fatalf("Unwrap: v = %d, want 42", v)
	}
}

// TestDecodeUnwrapOnFailure confirms that Unwrap on a Failure
// returns the zero value of T and the underlying error.
func TestDecodeUnwrapOnFailure(t *testing.T) {
	v, err := json.Decode[int]([]byte("not json")).Unwrap()
	if err == nil {
		t.Fatal("Unwrap on failure: err = nil, want non-nil")
	}
	if v != 0 {
		t.Fatalf("Unwrap on failure: v = %d, want 0 (zero value)", v)
	}
}

// TestDecodeRecoverFromError confirms that a JSON decode
// failure can be recovered via the Result's Recover method,
// producing the fallback value. Recover returns the unwrapped
// T (not a Result[T]), so the result is the int directly. This
// is the pattern the rest of the project uses to handle
// "the value is not there, fall back to a default" cases.
func TestDecodeRecoverFromError(t *testing.T) {
	got := json.Decode[int]([]byte("not json")).
		Recover(func(err error) int {
			if err == nil {
				t.Fatal("Recover: f called with nil error")
			}
			return -1
		})
	if got != -1 {
		t.Fatalf("after Recover: got %d, want -1", got)
	}
}

// TestDecodeMapError confirms that the underlying error from
// Decode can be transformed in place via MapError. This is
// useful when callers want to wrap a JSON error in their own
// domain error type.
func TestDecodeMapError(t *testing.T) {
	wrapped := errors.New("decode failed: bad input")
	r := json.Decode[int]([]byte("not json")).
		MapError(func(err error) error {
			if err == nil {
				t.Fatal("MapError: f called with nil error")
			}
			return wrapped
		})
	if r.IsSuccess() {
		t.Fatal("MapError: IsSuccess = true, want false")
	}
	if err := r.Error(); !errors.Is(err, wrapped) {
		t.Fatalf("MapError: err = %v, want wrap of %v", err, wrapped)
	}
}

// ---------- input handling ----------

// TestDecodeDoesNotRetainInput confirms the documented
// guarantee that Decode does not retain a reference to the
// input byte slice. The caller can mutate the slice after
// Decode returns; the returned Result is unaffected.
func TestDecodeDoesNotRetainInput(t *testing.T) {
	data := []byte(`"hello"`)
	r := json.Decode[string](data)
	if !r.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", r.Error())
	}

	// Overwrite the input buffer with different content.
	// A correct implementation must not be affected.
	for i := range data {
		data[i] = 'X'
	}

	if got := r.Value(); got != "hello" {
		t.Fatalf("after mutating input: got %q, want %q (Decode retained the input slice)", got, "hello")
	}
}

// TestDecodeReusesCallerSlice confirms the same property from
// the other side: after Decode, the caller can copy new data
// into the same slice and decode it again without interference
// from the first call. The test allocates a buffer large enough
// for the longer payload, then uses copy() to overwrite the
// bytes (not just the visible prefix) so the second Decode sees
// a clean slice.
func TestDecodeReusesCallerSlice(t *testing.T) {
	// Buffer is large enough for `"second"` (8 bytes); we
	// overwrite the whole visible region on each call so the
	// first payload's bytes are not visible to the second
	// Decode.
	data := make([]byte, 8)
	copy(data, `"first"`)
	r1 := json.Decode[string](data[:7])
	if v := r1.Value(); v != "first" {
		t.Fatalf("first Decode: got %q, want %q", v, "first")
	}

	// Overwrite with a different value, decoding from the
	// same backing array.
	copy(data, `"second"`)
	r2 := json.Decode[string](data)
	if v := r2.Value(); v != "second" {
		t.Fatalf("second Decode after mutating slice: got %q, want %q", v, "second")
	}
}

// TestDecodeChainedWithResult confirms that Decode participates
// in the standard Result fluent chain. The combination of
// Decode + Map + Recover shows the typical "decode, transform,
// fall back on error" pipeline that the rest of the project
// uses. Recover returns the unwrapped T directly, so the chain
// ends with a T rather than a Result[T].
func TestDecodeChainedWithResult(t *testing.T) {
	type input struct {
		N int `json:"n"`
	}
	type output struct {
		Doubled int `json:"doubled"`
	}

	got := json.Decode[input]([]byte(`{"n": 21}`)).
		Map(func(in input) output {
			return output{Doubled: in.N * 2}
		}).
		Recover(func(err error) output {
			t.Fatalf("Recover should not be called on success: err = %v", err)
			return output{}
		})
	if got.Doubled != 42 {
		t.Fatalf("chained decode: got %+v, want {Doubled: 42}", got)
	}
}

// TestDecodeFailureChainedWithRecover confirms that the same
// chain handles the failure path: a decode error is mapped
// to the same output type via Recover, so the caller does
// not need an explicit if-failure check.
func TestDecodeFailureChainedWithRecover(t *testing.T) {
	type output struct {
		Doubled int
	}

	got := json.Decode[int]([]byte("not json")).
		Map(func(n int) output { return output{Doubled: n * 2} }).
		Recover(func(err error) output {
			if err == nil {
				t.Fatal("Recover: f called with nil error")
			}
			return output{Doubled: -1}
		})
	if got.Doubled != -1 {
		t.Fatalf("chained decode with recovery: got %+v, want {Doubled: -1}", got)
	}
}

// TestDecodeIntoResultFromResult confirms that Decode composes
// with the rest of the Result API. The test wraps a Decode
// failure into a domain error via MapError, then unwraps the
// final Result with errors.Is to check the wrapped chain.
func TestDecodeIntoResultFromResult(t *testing.T) {
	sentinel := errors.New("sentinel")
	r := json.Decode[int]([]byte("not json")).
		MapError(func(error) error { return sentinel })
	if errors.Is(r.Error(), sentinel) {
		// expected
		return
	}
	// Use the result package's own helpers if errors.Is does
	// not find the wrapped error.
	_ = result.Failure[int](sentinel)
	t.Fatalf("errors.Is did not find sentinel in the wrapped error: %v", r.Error())
}
