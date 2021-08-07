package crawler

import "time"

type Repository interface {
	GetEvents(from time.Time) ([]Event, error)
	SaveEvents(events []Event) error
	GetRestriction(time time.Time) (Restriction, error)
	SaveRestriction(restriction Restriction) error
}
