package crawler

import "time"

type Repository interface {
	GetEvents(from time.Time) ([]Event, error)
	SaveEvents(events []Event) error
	GetRestriction(crawler string, job string) (Restriction, error)
	SaveRestriction(restriction Restriction) error
	GetUser(userName string) (User, error)
	SaveUsers(users []User) error
}
