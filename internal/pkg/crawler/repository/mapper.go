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
	}
}

func (m mapper) mapEventsToModelEvents(events []crawler.Event) []Event {
	var modelEvents []Event
	for _, e := range events {
		modelEvents = append(modelEvents, m.mapEventToModelEvent(e))
	}
	return modelEvents
}

func (m mapper) mapModelEventToEvent(e Event) crawler.Event {
	return crawler.Event{
		Crawler:  e.Crawler,
		Job:      e.Job,
		UserName: e.UserName,
		ID:       e.ID,
	}
}

func (m mapper) mapModelEventsToEvents(events []Event) []crawler.Event {
	var modelEvents []crawler.Event
	for _, e := range events {
		modelEvents = append(modelEvents, m.mapModelEventToEvent(e))
	}
	return modelEvents
}

func (m mapper) mapUserToModelUser(u crawler.User) User {
	return User{
		ID:   u.ID,
		Name: u.Name,
	}
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
