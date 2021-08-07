package crawler

type InvalidUserNameError struct {
	message string
}

func (e InvalidUserNameError) Error() string {
	return e.message
}
