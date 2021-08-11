package crawler

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

type UseCase struct {
	logger      *zap.Logger
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

	events, err := u.filterEvents(crawledEvents)
	if err != nil {
		return err
	}

	return u.notify(events)
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

// save events and returns saved events
func (u UseCase) filterEvents(crawledEvents []Event) ([]Event, error) {
	var events []Event
	for _, e := range crawledEvents {
		err := u.repository.SaveEvent(e)
		if errors.As(err, &AlreadyExistsError{}) {
			continue
		} else if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
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
				"notify error",
				zap.Error(err),
				zap.Int("index", i),
				zap.Any("event", e),
				zap.Any("events", events),
				zap.Any("notification", n),
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

		user, err = u.repository.GetUser(userName)
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
	logger *zap.Logger,
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
