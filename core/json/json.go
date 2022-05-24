package json

import (
	"encoding/json"
	"github.com/qianwj/typed/core/result"
	"io"
)

func Marshal[T any](data T) result.Result[[]byte] {
	bytes, err := json.Marshal(data)
	return result.Wrap[[]byte](bytes, err)
}

func Unmarshal[T any](data []byte) result.Result[T] {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return result.Err[T](err)
	}
	return result.Ok[T](t)
}

func Encode[T any](w io.Writer, data T) error {
	return json.NewEncoder(w).Encode(data)
}

func Decode[T any](r io.Reader) result.Result[T] {
	var t T
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return result.Err[T](err)
	}
	return result.Ok[T](t)
}
