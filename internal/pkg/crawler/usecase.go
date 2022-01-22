package crawler

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

type UseCase struct {
	logger         *zap.Logger
	repository     Repository
	crawler        Crawler
	notifier       Notifier
	channelService ChannelService
}

// TODO rename
func (u UseCase) Work(crawler string, job string) error {
	allowed, err := u.isAllowed(crawler, job)
	if err != nil {
		return err
	} else if !allowed {
		u.logger.Info("not allowed to run command")
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

func (u UseCase) isAllowed(crawler string, job string) (bool, error) {
	r, err := u.repository.GetRestriction(crawler, job)
	if err != nil {
		return false, err
	}

	now := time.Now().Add(time.Hour * 9) // kst변환
	u.logger.Info("isAllowed()", zap.Any("restriction", r))

	// no restriction record
	if now.After(r.StartDate) && now.Before(r.EndDate) {
		return true, nil
	}

	// allow cond   return
	// true  true   true	(run command)
	// true  false  false
	// false true   false
	// false false  true	(run commmand)
	cond := r.HourFrom <= now.Hour() && now.Hour() < r.HourTo
	if r.Allow {
		if cond {
			return true, nil
		} else { // !cond
			return false, nil
		}
	} else { // !r.Allow
		if cond {
			return false, nil
		} else { // !cond
			return true, nil
		}
	}
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

func (u UseCase) getUser(userName string) (Channel, error) {
	user, err := u.repository.GetUser(userName)
	if err != nil {
		return Channel{}, err
	} else if user.ID == "" {
		// sync with slack
		users, err := u.channelService.GetChannels()
		if err != nil {
			return Channel{}, err
		}

		err = u.repository.SaveUsers(users)
		if err != nil {
			return Channel{}, err
		}

		user, err = u.repository.GetUser(userName)
		if err != nil {
			return Channel{}, err
		} else if user.ID == "" {
			return Channel{}, errors.New("empty user")
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
	channelService ChannelService,
) *UseCase {
	return &UseCase{
		logger:         logger,
		repository:     repository,
		crawler:        crawler,
		notifier:       notifier,
		channelService: channelService,
	}
}
