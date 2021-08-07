package crawler

import (
	"time"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/logger"
)

type UseCase struct {
	repository Repository
	crawler    Crawler
	notifier   Notifier
	logger     logger.Logger
}

// TODO rename
func (u UseCase) Work() error {
	now := time.Now()

	restricted, err := u.isRestricted(now)
	if err != nil {
		return err
	} else if restricted {
		u.logger.Info("restricted")
		return nil
	}

	crawledEvents, err := u.crawler.Crawl()
	if err != nil {
		return err
	}

	notifiedEvents, err := u.repository.GetEvents(now.AddDate(0, 0, -1))
	if err != nil {
		return err
	}
	newEvents := u.filterEvents(crawledEvents, notifiedEvents)

	err = u.repository.SaveEvents(newEvents)
	if err != nil {
		return err
	}

	return u.notify(newEvents)
}

func (u UseCase) isRestricted(t time.Time) (bool, error) {
	restriction, err := u.repository.GetRestriction(t)
	if err != nil {
		return false, err
	}

	return restriction.Crawler != "", nil
}

func (u UseCase) filterEvents(crawledEvents []Event, notifiedEvents []Event) []Event {
	toMap := func(events []Event) map[string]Event {
		m := map[string]Event{}
		for _, e := range events {
			m[e.ID] = e
		}
		return m
	}

	crawled := toMap(crawledEvents)
	notified := toMap(notifiedEvents)
	new := []Event{}
	for id, e := range crawled {
		_, ok := notified[id]
		if !ok {
			new = append(new, e)
		}
	}

	return new
}

func (u UseCase) notify(events []Event) error {
	for i, e := range events {
		err := u.notifier.Notify(e)
		if err != nil {
			u.logger.Error(
				"error on notify",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "index", Value: i},
				logger.Field{Key: "event", Value: e},
				logger.Field{Key: "events", Value: events},
			)
			return err
		}
	}

	return nil
}

func (u UseCase) AddRestriction(restriction Restriction) error {
	return u.repository.SaveRestriction(restriction)
}

func NewUseCase(
	logger logger.Logger,
	repository Repository,
	crawler Crawler,
	notifier Notifier,
) *UseCase {
	return &UseCase{
		logger:     logger,
		repository: repository,
		crawler:    crawler,
		notifier:   notifier,
	}
}
