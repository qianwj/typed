package errors

var (
	errorConfNotFound = &confNotFoundError{}
)

type confNotFoundError struct{}

func ConfNotFound() error {
	return errorConfNotFound
}

func (e *confNotFoundError) Error() string {
	return "configuration not found"
}
