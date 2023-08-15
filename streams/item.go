package streams

type Item struct {
	Data  any
	Error error
}

func Ok(data any) *Item {
	return &Item{Data: data}
}

func Err(err error) *Item {
	return &Item{Error: err}
}
