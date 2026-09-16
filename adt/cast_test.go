package adt_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/qianwj/typed/adt"
)

func TestCastSuccess(t *testing.T) {
	for _, want := range []int{0, 42} {
		got := adt.Cast[int](want)
		if got.IsFailure() || got.Value() != want {
			t.Fatalf("Cast[int](%d) = %v", want, got)
		}
	}
	buf := new(bytes.Buffer)
	got := adt.Cast[io.Reader](buf)
	if got.IsFailure() || got.Value() != buf {
		t.Fatalf("Cast[io.Reader] lost implementing value: %v", got)
	}
}

func TestCastFailure(t *testing.T) {
	type count int
	for _, tc := range []struct {
		name       string
		value      any
		actualType string
	}{
		{"wrong type", "42", "string"},
		{"numeric conversion", int32(42), "int32"},
		{"distinct named type", count(42), "adt_test.count"},
		{"nil interface", nil, "<nil>"},
		{"wrong typed nil", (*int)(nil), "*int"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := adt.Cast[int](tc.value)
			value, err := got.Unwrap()
			if !got.IsFailure() || err == nil || value != 0 {
				t.Fatalf("Cast[int](%v) = %v, want failure", tc.value, got)
			}
			if !strings.Contains(err.Error(), tc.actualType) || !strings.Contains(err.Error(), "as int") {
				t.Fatalf("error missing source or target type: %v", err)
			}
		})
	}
	got := adt.Cast[io.Reader](42)
	if got.IsSuccess() || !strings.Contains(got.Error().Error(), "io.Reader") {
		t.Fatalf("interface target missing from failure: %v", got)
	}
	if got := adt.Cast[any](nil); got.IsSuccess() {
		t.Fatal("nil interface must fail even when T is any")
	}
}

func TestCastTypedNilSuccess(t *testing.T) {
	var pointer *bytes.Buffer
	if got := adt.Cast[*bytes.Buffer](pointer); got.IsFailure() || got.Value() != nil {
		t.Fatalf("typed nil pointer should succeed: %v", got)
	}
	if got := adt.Cast[io.Reader](pointer); got.IsFailure() || got.Value() != pointer {
		t.Fatalf("typed nil implementing an interface should succeed: %v", got)
	}
	if got := adt.Cast[any](pointer); got.IsFailure() || got.Value() != pointer {
		t.Fatalf("typed nil dynamic type should be retained: %v", got)
	}
	var values []int
	if got := adt.Cast[[]int](values); got.IsFailure() || got.Value() != nil {
		t.Fatalf("typed nil slice should succeed: %v", got)
	}
}
