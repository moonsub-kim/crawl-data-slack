package crawler

type Repository interface {
	SaveEvent(event Event) error
	GetRestriction(crawler string, job string) (Restriction, error)
	SaveRestriction(restriction Restriction) error
	GetUser(userName string) (User, error)
	SaveUsers(users []User) error
}
