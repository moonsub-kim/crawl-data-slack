package crawler

import (
	"errors"
	"time"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/logger"
)

type UseCase struct {
	logger      logger.Logger
	repository  Repository
	crawler     Crawler
	notifier    Notifier
	userService UserService
}

// TODO rename
func (u UseCase) Work(crawler string, job string) error {
	restricted, err := u.isRestricted(crawler, job)
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

	notifiedEvents, err := u.repository.GetEvents(time.Now().AddDate(0, 0, -1))
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

func (u UseCase) isRestricted(crawler string, job string) (bool, error) {
	r, err := u.repository.GetRestriction(crawler, job)
	if err != nil {
		return false, err
	}

	now := time.Now()
	if now.After(r.StartDate) && now.Before(r.EndDate) && r.HourFrom <= now.Hour() && now.Hour() < r.HourTo {
		return true, nil
	}

	return false, nil
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
		user, err := u.getUser(e.UserName)
		if err != nil {
			return err
		}

		n := Notification{
			Event: e,
			User:  user,
		}

		err = u.notifier.Notify(n)
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

func (u UseCase) getUser(userName string) (User, error) {
	user, err := u.repository.GetUser(userName)
	if err != nil {
		return User{}, err
	} else if user.ID == "" {
		// sync with slack
		users, err := u.userService.GetUsers()
		if err != nil {
			return User{}, err
		}

		err = u.repository.SaveUsers(users)
		if err != nil {
			return User{}, err
		}

		user, err := u.repository.GetUser(userName)
		if err != nil {
			return User{}, err
		} else if user.ID == "" {
			return User{}, errors.New("empty user")
		}
	}

	return user, nil
}

func (u UseCase) AddRestriction(restriction Restriction) error {
	return u.repository.SaveRestriction(restriction)
}

func NewUseCase(
	logger logger.Logger,
	repository Repository,
	crawler Crawler,
	notifier Notifier,
	userService UserService,
) *UseCase {
	return &UseCase{
		logger:      logger,
		repository:  repository,
		crawler:     crawler,
		notifier:    notifier,
		userService: userService,
	}
}
