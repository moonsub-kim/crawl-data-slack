package crawler

import (
	"time"
)

type Event struct {
	Crawler  string
	Job      string
	UserName string // <firstname>.<lastname>
	UID      string // UID is determined by Crawler
	Name     string
	Message  string
}

type User struct {
	ID   string
	Name string
}

type Restriction struct {
	Crawler   string
	Job       string
	StartDate time.Time
	EndDate   time.Time
	HourFrom  int
	HourTo    int
}

type Notification struct {
	Event Event
	User  User
}
