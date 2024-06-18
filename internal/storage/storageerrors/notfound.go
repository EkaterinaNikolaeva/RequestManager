package storageerrors

type NotFoundError struct {
}

func NewNotFoundError() NotFoundError {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "the object was not found in the storage"
}
