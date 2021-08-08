package repository

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
)

type mapper struct {
}

func (m mapper) mapEventToModelEvent(e crawler.Event) Event {
	return Event{
		Crawler:  e.Crawler,
		Job:      e.Job,
		UserName: e.UserName,
		ID:       e.ID,
		Name:     e.Name,
		Message:  e.Message,
	}
}

func (m mapper) mapUserToModelUser(u crawler.User) User {
	return User{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (m mapper) mapUsersToModelUsers(users []crawler.User) []User {
	var modelUsers []User
	for _, u := range users {
		modelUsers = append(modelUsers, m.mapUserToModelUser(u))
	}
	return modelUsers
}

func (m mapper) mapModelUserToUser(u User) crawler.User {
	return crawler.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (m mapper) mapModelRestrictionToRestriction(r Restriction) crawler.Restriction {
	return crawler.Restriction{
		Crawler:   r.Crawler,
		Job:       r.Job,
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
		HourFrom:  r.HourFrom,
		HourTo:    r.HourTo,
	}
}

func (m mapper) mapRestrictionToModelRestriction(r crawler.Restriction) Restriction {
	return Restriction{
		Crawler:   r.Crawler,
		Job:       r.Job,
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
		HourFrom:  r.HourFrom,
		HourTo:    r.HourTo,
	}
}
