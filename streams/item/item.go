package item

func Ok(data any) Item {
	return Item{data: data}
}

func Err(err error) Item {
	return Item{error: err}
}

type Item struct {
	data  any
	error error
}

func (i Item) IsOk() bool {
	return i.error == nil
}

func (i Item) IsErr() bool {
	return i.error != nil
}

func (i Item) Data() any {
	return i.data
}

func (i Item) Err() error {
	return i.error
}

func (i Item) Get() (any, error) {
	return i.data, i.error
}
