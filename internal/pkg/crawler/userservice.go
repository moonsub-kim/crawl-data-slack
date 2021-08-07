package crawler

type UserService interface {
	GetUsers() ([]User, error)
}
