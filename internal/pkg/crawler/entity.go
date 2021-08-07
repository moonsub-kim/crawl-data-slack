package crawler

import (
	"time"
)

type Event struct {
	Crawler  string
	Job      string
	UserName string
	ID       string // ID is determined by Crawler
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
	HourFrom  time.Time
	HourTo    time.Time
}
