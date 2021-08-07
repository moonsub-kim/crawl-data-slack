package crawler

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	Crawler  string
	Name     string
	UserName UserName
	ID       string // ID is determined by Crawler
}

type User struct {
	ID   string
	Name UserName
}

type UserName struct {
	FirstName string
	LastName  string
}

func (u UserName) String() string {
	return u.FirstName + " " + u.LastName
}

func (u UserName) Assert() error {
	if strings.ToLower(u.FirstName) != u.FirstName || strings.ToLower(u.LastName) != u.LastName {
		return InvalidUserNameError{message: fmt.Sprintf("UserName \"%s\" has upper case letter.", u)}
	}
	return nil
}

type Restriction struct {
	ID        int
	StartDate time.Time
	EndDate   time.Time
	HourFrom  time.Time
	HourTo    time.Time
}
