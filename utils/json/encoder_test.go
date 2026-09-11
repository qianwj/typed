package json_test

import (
	v2 "encoding/json/v2"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/qianwj/typed/utils/json"
	"github.com/qianwj/typed/utils/result"
)

// ---------- success path: primitive values ----------

// TestEncodeInteger confirms the basic happy path for a
// primitive numeric value.
func TestEncodeInteger(t *testing.T) {
	r := json.Encode(42)
	if !r.IsSuccess() {
		t.Fatalf("Encode(42): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "42" {
		t.Fatalf("Encode(42): got %s, want 42", got)
	}
}

// TestEncodeString confirms the happy path for a string.
// The output must include the surrounding quotes.
func TestEncodeString(t *testing.T) {
	r := json.Encode("hello")
	if !r.IsSuccess() {
		t.Fatalf("Encode(\"hello\"): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != `"hello"` {
		t.Fatalf("Encode(\"hello\"): got %s, want %q", got, `"hello"`)
	}
}

// TestEncodeBool confirms the happy path for a boolean.
func TestEncodeBool(t *testing.T) {
	r := json.Encode(true)
	if !r.IsSuccess() {
		t.Fatalf("Encode(true): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "true" {
		t.Fatalf("Encode(true): got %s, want true", got)
	}
}

// TestEncodeFloat confirms the happy path for a float value
// that is not an integer.
func TestEncodeFloat(t *testing.T) {
	r := json.Encode(3.14)
	if !r.IsSuccess() {
		t.Fatalf("Encode(3.14): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); !strings.HasPrefix(string(got), "3.14") {
		t.Fatalf("Encode(3.14): got %s, want 3.14 prefix", got)
	}
}

// ---------- success path: composite values ----------

// TestEncodeStruct confirms the happy path for a struct
// value, with field names matching the exported field
// names.
func TestEncodeStruct(t *testing.T) {
	type profile struct {
		Name string
		Age  int
	}
	r := json.Encode(profile{Name: "alice", Age: 30})
	if !r.IsSuccess() {
		t.Fatalf("Encode(profile): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != `{"Name":"alice","Age":30}` {
		t.Fatalf("Encode(profile): got %s, want %q", got, `{"Name":"alice","Age":30}`)
	}
}

// TestEncodeStructWithTags confirms that struct tags are
// honoured on the encode path. The output uses the tag
// names ("first_name", "last_name") rather than the field
// names.
func TestEncodeStructWithTags(t *testing.T) {
	type profile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	r := json.Encode(profile{FirstName: "alice", LastName: "smith"})
	if !r.IsSuccess() {
		t.Fatalf("Encode(profile): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != `{"first_name":"alice","last_name":"smith"}` {
		t.Fatalf("Encode(profile): got %s, want %q", got, `{"first_name":"alice","last_name":"smith"}`)
	}
}

// TestEncodeSlice confirms the happy path for a slice
// value.
func TestEncodeSlice(t *testing.T) {
	r := json.Encode([]int{1, 2, 3, 4, 5})
	if !r.IsSuccess() {
		t.Fatalf("Encode(slice): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "[1,2,3,4,5]" {
		t.Fatalf("Encode(slice): got %s, want [1,2,3,4,5]", got)
	}
}

// TestEncodeMap confirms the happy path for a string-keyed
// map value.
func TestEncodeMap(t *testing.T) {
	r := json.Encode(map[string]int{"a": 1, "b": 2, "c": 3})
	if !r.IsSuccess() {
		t.Fatalf("Encode(map): IsSuccess = false, err = %v", r.Error())
	}
	// Map iteration order is non-deterministic, so the
	// test parses the output and compares the entries
	// as a set.
	var got map[string]int
	if err := v2.Unmarshal(r.Value(), &got); err != nil {
		t.Fatalf("parse Encode output: %v", err)
	}
	want := map[string]int{"a": 1, "b": 2, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Encode(map): got %v, want %v", got, want)
	}
}

// TestEncodeNested confirms the happy path for a nested
// struct.
func TestEncodeNested(t *testing.T) {
	type address struct {
		City string
		Zip  string
	}
	type person struct {
		Name    string
		Address address
	}
	r := json.Encode(person{Name: "alice", Address: address{City: "Boston", Zip: "02115"}})
	if !r.IsSuccess() {
		t.Fatalf("Encode(person): IsSuccess = false, err = %v", r.Error())
	}
	want := `{"Name":"alice","Address":{"City":"Boston","Zip":"02115"}}`
	if got := r.Value(); string(got) != want {
		t.Fatalf("Encode(person): got %s, want %s", got, want)
	}
}

// TestEncodePointer confirms that a pointer to a struct
// marshals the struct (not the pointer's address).
func TestEncodePointer(t *testing.T) {
	type inner struct{ V int }
	v := &inner{V: 7}
	r := json.Encode(v)
	if !r.IsSuccess() {
		t.Fatalf("Encode(pointer): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != `{"V":7}` {
		t.Fatalf("Encode(pointer): got %s, want %q", got, `{"V":7}`)
	}
}

// TestEncodeNilPointer confirms that a nil pointer to a
// struct marshals to "null".
func TestEncodeNilPointer(t *testing.T) {
	type inner struct{ V int }
	var p *inner
	r := json.Encode(p)
	if !r.IsSuccess() {
		t.Fatalf("Encode(nil pointer): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "null" {
		t.Fatalf("Encode(nil pointer): got %s, want null", got)
	}
}

// TestEncodeNilSlice confirms that a nil slice marshals to
// "[]" by default in encoding/json/v2. This is a v2
// behavioural difference from v1: the v1 json package
// rendered a nil slice as "null" by default; v2 renders it
// as "[]". A separate test (TestEncodeNilSliceAsNullOption)
// confirms that the FormatNilSliceAsNull option flips the
// behaviour to the v1 default.
func TestEncodeNilSlice(t *testing.T) {
	var s []int
	r := json.Encode(s)
	if !r.IsSuccess() {
		t.Fatalf("Encode(nil slice): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "[]" {
		t.Fatalf("Encode(nil slice): got %s, want [] (v2 default)", got)
	}
}

// TestEncodeNilSliceAsNullOption confirms that the
// FormatNilSliceAsNull option produces "null" for a nil
// slice, matching v1 behaviour.
func TestEncodeNilSliceAsNullOption(t *testing.T) {
	var s []int
	r := json.Encode(s, v2.FormatNilSliceAsNull(true))
	if !r.IsSuccess() {
		t.Fatalf("Encode(nil slice, option): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "null" {
		t.Fatalf("Encode(nil slice, option): got %s, want null (with FormatNilSliceAsNull)", got)
	}
}

// TestEncodeEmptySliceIsArray confirms that a non-nil but
// empty slice marshals to "[]" (not "null"). This is the
// distinction the doc draws between nil and empty.
func TestEncodeEmptySliceIsArray(t *testing.T) {
	s := []int{}
	r := json.Encode(s)
	if !r.IsSuccess() {
		t.Fatalf("Encode(empty slice): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != "[]" {
		t.Fatalf("Encode(empty slice): got %s, want []", got)
	}
}

// ---------- success path: custom Marshaler ----------

// TestEncodeCustomMarshaler confirms that a type that
// implements json.Marshaler is encoded by its MarshalJSON
// method. The custom encoding completely replaces the
// default struct encoding.
func TestEncodeCustomMarshaler(t *testing.T) {
	type point struct {
		X, Y int
	}

	// Define a wrapper type with a custom MarshalJSON that
	// always emits "P(X,Y)" so we can verify the method
	// is being called.
	type custom struct{ p point }
	_ = custom{}

	// Use a value type that has its own MarshalJSON.
	r := json.Encode(marshalerPoint{X: 1, Y: 2})
	if !r.IsSuccess() {
		t.Fatalf("Encode(marshaler): IsSuccess = false, err = %v", r.Error())
	}
	if got := r.Value(); string(got) != `"P(1,2)"` {
		t.Fatalf("Encode(marshaler): got %s, want %q", got, `"P(1,2)"`)
	}
}

// marshalerPoint is a test type that implements json.Marshaler
// by always emitting a fixed string. This lets the test
// confirm that the Marshaler interface is honoured by Encode.
type marshalerPoint struct{ X, Y int }

func (m marshalerPoint) MarshalJSON() ([]byte, error) {
	return []byte(`"P(` + itoa(m.X) + `,` + itoa(m.Y) + `)"`), nil
}

// itoa is a minimal int-to-string converter that avoids the
// strconv import in the test file.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// ---------- options ----------

// TestEncodeDeterministicMap confirms that the
// Deterministic option produces a stable output for a map.
// Without the option, map iteration order is not guaranteed;
// with it, the keys are emitted in sorted order.
func TestEncodeDeterministicMap(t *testing.T) {
	src := map[string]int{"c": 3, "a": 1, "b": 2}
	r1 := json.Encode(src, v2.Deterministic(true))
	r2 := json.Encode(src, v2.Deterministic(true))
	if !r1.IsSuccess() || !r2.IsSuccess() {
		t.Fatalf("Encode(deterministic): failures: %v, %v", r1.Error(), r2.Error())
	}
	if string(r1.Value()) != string(r2.Value()) {
		t.Fatalf("Encode(deterministic): non-stable: %s vs %s", r1.Value(), r2.Value())
	}
	// Sorted order must be a, b, c.
	if got := string(r1.Value()); got != `{"a":1,"b":2,"c":3}` {
		t.Fatalf("Encode(deterministic): got %s, want %q", got, `{"a":1,"b":2,"c":3}`)
	}
}

// TestEncodeOmitZeroStructFields confirms that the
// OmitZeroStructFields option drops fields whose value is
// the type's zero value.
func TestEncodeOmitZeroStructFields(t *testing.T) {
	type profile struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	r := json.Encode(profile{Name: "alice"},
		v2.Deterministic(true),
		v2.OmitZeroStructFields(true))
	if !r.IsSuccess() {
		t.Fatalf("Encode(omit): IsSuccess = false, err = %v", r.Error())
	}
	if got := string(r.Value()); got != `{"name":"alice"}` {
		t.Fatalf("Encode(omit): got %s, want %q", got, `{"name":"alice"}`)
	}
}

// TestEncodeWithoutOptions confirms that calling Encode
// without any options is valid and produces the same output
// as the default v2 marshaling.
func TestEncodeWithoutOptions(t *testing.T) {
	r := json.Encode([]int{1, 2, 3})
	if !r.IsSuccess() {
		t.Fatalf("Encode(slice): IsSuccess = false, err = %v", r.Error())
	}
	if got := string(r.Value()); got != "[1,2,3]" {
		t.Fatalf("Encode(slice): got %s, want [1,2,3]", got)
	}
}

// ---------- failure path ----------

// TestEncodeUnencodableFunction confirms that a value
// containing a Go function (which cannot be serialised)
// produces a Failure with a non-nil error.
func TestEncodeUnencodableFunction(t *testing.T) {
	r := json.Encode(struct {
		F func()
	}{F: func() {}})
	if r.IsSuccess() {
		t.Fatal("Encode(function): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Encode(function): Error = nil, want non-nil")
	}
}

// TestEncodeUnencodableChannel confirms that a value
// containing a Go channel produces a Failure.
func TestEncodeUnencodableChannel(t *testing.T) {
	r := json.Encode(struct {
		C chan int
	}{C: make(chan int)})
	if r.IsSuccess() {
		t.Fatal("Encode(channel): IsSuccess = true, want false")
	}
	if err := r.Error(); err == nil {
		t.Fatal("Encode(channel): Error = nil, want non-nil")
	}
}

// TestEncodeMarshalerError confirms that a custom
// Marshaler returning a non-nil error is propagated as a
// Failure. v2 wraps the original error in a "json: cannot
// marshal from Go ..." prefix, so the test uses
// strings.Contains rather than exact-match to verify the
// original message is preserved.
type erroringMarshaler struct{}

func (erroringMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshaler failed")
}

func TestEncodeMarshalerError(t *testing.T) {
	r := json.Encode(erroringMarshaler{})
	if r.IsSuccess() {
		t.Fatal("Encode(marshaler error): IsSuccess = true, want false")
	}
	err := r.Error()
	if err == nil {
		t.Fatal("Encode(marshaler error): Error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "marshaler failed") {
		t.Fatalf("Encode(marshaler error): err = %v, want message containing %q", err, "marshaler failed")
	}
}

// ---------- Unwrap and bridge ----------

// TestEncodeUnwrapToStandardShape confirms the recommended
// pattern: Unwrap extracts the standard ([]byte, error)
// pair from a Success.
func TestEncodeUnwrapToStandardShape(t *testing.T) {
	data, err := json.Encode(42).Unwrap()
	if err != nil {
		t.Fatalf("Unwrap: err = %v, want nil", err)
	}
	if string(data) != "42" {
		t.Fatalf("Unwrap: data = %s, want 42", data)
	}
}

// TestEncodeUnwrapOnFailure confirms that Unwrap on a
// Failure returns a nil []byte and the underlying error.
func TestEncodeUnwrapOnFailure(t *testing.T) {
	data, err := json.Encode(struct {
		F func()
	}{F: func() {}}).Unwrap()
	if err == nil {
		t.Fatal("Unwrap on failure: err = nil, want non-nil")
	}
	if data != nil {
		t.Fatalf("Unwrap on failure: data = %s, want nil", data)
	}
}

// TestEncodeOrElseFailure confirms the OrElse fallback
// pattern. On a Failure, OrElse returns the caller-supplied
// fallback bytes.
func TestEncodeOrElseFailure(t *testing.T) {
	r := json.Encode(struct {
		F func()
	}{F: func() {}}).
		OrElse([]byte(`{"default":true}`))
	if string(r) != `{"default":true}` {
		t.Fatalf("OrElse: got %s, want %q", r, `{"default":true}`)
	}
}

// TestEncodeOrElseSuccess confirms that OrElse on a
// Success returns the encoded bytes, not the fallback.
func TestEncodeOrElseSuccess(t *testing.T) {
	r := json.Encode(42).OrElse([]byte("99"))
	if string(r) != "42" {
		t.Fatalf("OrElse on success: got %s, want 42", r)
	}
}

// TestEncodeRecover confirms the Recover pattern: a
// failure from Encode is turned into a Success with a
// fallback value (the encoded bytes).
func TestEncodeRecover(t *testing.T) {
	r := json.Encode(struct {
		F func()
	}{F: func() {}}).
		Recover(func(err error) []byte {
			if err == nil {
				t.Fatal("Recover: f called with nil error")
			}
			return []byte(`{"recovered":true}`)
		})
	if string(r) != `{"recovered":true}` {
		t.Fatalf("Recover: got %s, want %q", r, `{"recovered":true}`)
	}
}

// TestEncodeMapError confirms that the underlying error
// can be transformed in place via MapError. This is useful
// when callers want to wrap an Encode error in their own
// domain error type.
func TestEncodeMapError(t *testing.T) {
	wrapped := errors.New("encode failed: bad value")
	r := json.Encode(struct {
		F func()
	}{F: func() {}}).
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

// ---------- round-trip with Decode ----------

// TestEncodeDecodeRoundTripPrimitive confirms that a value
// encoded by Encode can be decoded back by Decode into
// the same value.
func TestEncodeDecodeRoundTripPrimitive(t *testing.T) {
	cases := []int{0, 1, -1, 42, 1000, -1000}
	for _, want := range cases {
		encoded := json.Encode(want)
		if !encoded.IsSuccess() {
			t.Fatalf("Encode(%d): IsSuccess = false, err = %v", want, encoded.Error())
		}
		decodedResult := json.Decode[int](encoded.Value())
		if !decodedResult.IsSuccess() {
			t.Fatalf("Decode(Encode(%d)): IsSuccess = false, err = %v", want, decodedResult.Error())
		}
		if got := decodedResult.Value(); got != want {
			t.Fatalf("round-trip %d: got %d", want, got)
		}
	}
}

// TestEncodeDecodeRoundTripStruct confirms the round-trip
// property holds for struct values, including those with
// json tags.
func TestEncodeDecodeRoundTripStruct(t *testing.T) {
	type profile struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Tags    []string `json:"tags"`
		Friends []int    `json:"friends"`
	}
	src := profile{
		Name:    "alice",
		Age:     30,
		Tags:    []string{"admin", "user"},
		Friends: []int{1, 2, 3},
	}
	encoded := json.Encode(src)
	if !encoded.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", encoded.Error())
	}
	decoded := json.Decode[profile](encoded.Value())
	if !decoded.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", decoded.Error())
	}
	dst := decoded.Value()
	if !reflect.DeepEqual(dst, src) {
		t.Fatalf("round-trip struct: got %+v, want %+v", dst, src)
	}
}

// TestEncodeDecodeRoundTripSlice confirms the round-trip
// property for slices, where element order is preserved.
func TestEncodeDecodeRoundTripSlice(t *testing.T) {
	src := []int{3, 1, 4, 1, 5, 9, 2, 6}
	encoded := json.Encode(src)
	if !encoded.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", encoded.Error())
	}
	dst := json.Decode[[]int](encoded.Value())
	if !dst.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", dst.Error())
	}
	if !reflect.DeepEqual(dst.Value(), src) {
		t.Fatalf("round-trip slice: got %v, want %v", dst.Value(), src)
	}
}

// TestEncodeDecodeRoundTripMap confirms the round-trip
// property for maps. Map key sets must match (order is
// not preserved by Encode, so a set comparison is used).
func TestEncodeDecodeRoundTripMap(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2, "c": 3}
	encoded := json.Encode(src, v2.Deterministic(true))
	if !encoded.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", encoded.Error())
	}
	dst := json.Decode[map[string]int](encoded.Value())
	if !dst.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", dst.Error())
	}
	got := dst.Value()
	if len(got) != len(src) {
		t.Fatalf("round-trip map: got %d entries, want %d", len(got), len(src))
	}
	for k, v := range src {
		if gv, ok := got[k]; !ok || gv != v {
			t.Fatalf("round-trip map: key %q = %d, want %d (ok=%v)", k, gv, v, ok)
		}
	}
}

// TestEncodeDecodeRoundTripPointer confirms the round-trip
// property for pointer values. A non-nil pointer is
// preserved as a non-nil pointer on decode.
func TestEncodeDecodeRoundTripPointer(t *testing.T) {
	type inner struct{ V int }
	src := &inner{V: 42}
	encoded := json.Encode(src)
	if !encoded.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", encoded.Error())
	}
	dst := json.Decode[*inner](encoded.Value())
	if !dst.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", dst.Error())
	}
	got := dst.Value()
	if got == nil {
		t.Fatal("round-trip pointer: got nil, want non-nil")
	}
	if got.V != 42 {
		t.Fatalf("round-trip pointer: got V = %d, want 42", got.V)
	}
}

// TestEncodeDecodeRoundTripNested confirms the round-trip
// for nested composite types.
func TestEncodeDecodeRoundTripNested(t *testing.T) {
	type address struct {
		City string `json:"city"`
		Zip  string `json:"zip"`
	}
	type person struct {
		Name    string  `json:"name"`
		Address address `json:"address"`
	}
	src := person{
		Name:    "alice",
		Address: address{City: "Boston", Zip: "02115"},
	}
	encoded := json.Encode(src)
	if !encoded.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", encoded.Error())
	}
	dst := json.Decode[person](encoded.Value())
	if !dst.IsSuccess() {
		t.Fatalf("Decode: IsSuccess = false, err = %v", dst.Error())
	}
	if !reflect.DeepEqual(dst.Value(), src) {
		t.Fatalf("round-trip nested: got %+v, want %+v", dst.Value(), src)
	}
}

// ---------- chaining and composition ----------

// TestEncodeInResultChain confirms that Encode can sit at
// the start of a Result-based pipeline. The example builds
// a domain value, encodes it, then runs a post-processing
// step on the bytes (length check) wrapped in a Result.
func TestEncodeInResultChain(t *testing.T) {
	type profile struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	got := json.Encode(profile{Name: "alice", Age: 30}).
		Map(func(b []byte) int { return len(b) }).
		OrElse(0)
	if got == 0 {
		t.Fatalf("Map chain: got 0, want > 0 (encode succeeded)")
	}
}

// TestEncodeInResultChainOnFailure confirms the chain
// handles the failure path: a non-encodable value is
// turned into a default via Recover.
func TestEncodeInResultChainOnFailure(t *testing.T) {
	got := json.Encode(struct {
		F func()
	}{F: func() {}}).
		Map(func(b []byte) int { return len(b) }).
		Recover(func(err error) int {
			if err == nil {
				t.Fatal("Recover: f called with nil error")
			}
			return -1
		})
	if got != -1 {
		t.Fatalf("Recover chain: got %d, want -1 (fallback)", got)
	}
}

// TestEncodeReuseResult confirms that the Result returned
// by Encode can be inspected via IsSuccess / IsFailure
// without consuming the value.
func TestEncodeReuseResult(t *testing.T) {
	r := json.Encode("hello")
	if !r.IsSuccess() {
		t.Fatal("IsSuccess = false, want true")
	}
	if r.IsFailure() {
		t.Fatal("IsFailure = true, want false")
	}
	// Calling Value twice returns the same bytes.
	first := r.Value()
	second := r.Value()
	if string(first) != string(second) {
		t.Fatalf("Value twice: got %s and %s, want same", first, second)
	}
}

// ---------- sanity / regression ----------

// TestEncodeDeterministicSortKeys confirms that under
// Deterministic(true), map keys are sorted in lexicographic
// order. This is the property the doc relies on.
func TestEncodeDeterministicSortKeys(t *testing.T) {
	src := map[string]int{"z": 1, "m": 2, "a": 3, "b": 4}
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	r := json.Encode(src, v2.Deterministic(true))
	if !r.IsSuccess() {
		t.Fatalf("Encode: IsSuccess = false, err = %v", r.Error())
	}

	// Confirm the keys appear in sorted order in the output.
	output := string(r.Value())
	lastIdx := -1
	for _, k := range keys {
		idx := strings.Index(output, `"`+k+`"`)
		if idx == -1 {
			t.Fatalf("Encode(deterministic): key %q not in output %q", k, output)
		}
		if idx < lastIdx {
			t.Fatalf("Encode(deterministic): key %q out of order in %q", k, output)
		}
		lastIdx = idx
	}
}

// TestEncodeCompatibleWithResultImport confirms that the
// Encode return type composes with the rest of the
// result package. This is a smoke test of the type
// relationship: json.Encode returns result.Result[[]byte]
// and the result package methods work on it.
func TestEncodeCompatibleWithResultImport(t *testing.T) {
	// Construct an Encode result and exercise the result
	// API surface. This is a compile-time test: if the
	// types do not line up, the test fails to build.
	r := json.Encode(42)

	// Type assertion: r must be result.Result[[]byte].
	var _ result.Result[[]byte] = r

	// Use the API: Value, Error, IsSuccess, IsFailure,
	// Unwrap, OrElse, Optional.
	_ = r.Value()
	_ = r.Error()
	_ = r.IsSuccess()
	_ = r.IsFailure()
	_, _ = r.Unwrap()
	_ = r.OrElse(nil)
	_ = r.Optional()
}
