package crawler

type Repository interface {
	SaveEvent(event Event) error
	GetUser(userName string) (Channel, error)
	SaveUsers(users []Channel) error
}
