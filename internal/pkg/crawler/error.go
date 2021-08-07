package crawler

type AlreadyExistsError struct {
}

func (AlreadyExistsError) Error() string {
	return "already exists"
}
