package crawler

import (
	"errors"

	"go.uber.org/zap"
)

type UseCase struct {
	logger     *zap.Logger
	repository Repository
	crawler    Crawler
	messenger  Messenger
}

// TODO rename
func (u UseCase) Work() error {
	crawledEvents, err := u.crawler.Crawl()
	if err != nil {
		return err
	}

	u.logger.Info(
		"work",
		zap.Any("events", crawledEvents),
	)

	events, err := u.filterEvents(crawledEvents)
	if err != nil {
		return err
	}

	return u.notify(events)
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
		user, err := u.GetChannel(e.UserName)
		if err != nil {
			return err
		}

		n := Notification{
			Event: e,
			User:  user,
		}

		err = u.messenger.Notify(n)
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

func (u UseCase) GetChannel(name string) (Channel, error) {
	c, err := u.repository.GetChannel(name)
	if err != nil {
		return Channel{}, err
	} else if c.ID == "" {
		// sync with slack
		channels, err := u.messenger.GetChannels()
		if err != nil {
			return Channel{}, err
		}
		u.logger.Info("channels", zap.Any("channels", channels))

		err = u.repository.SyncChannels(channels)
		if err != nil {
			return Channel{}, err
		}

		c, err = u.repository.GetChannel(name)
		if err != nil {
			return Channel{}, err
		} else if c.ID == "" {
			return Channel{}, errors.New("empty channel")
		}
	}

	return c, nil
}

func NewUseCase(
	logger *zap.Logger,
	repository Repository,
	crawler Crawler,
	messenger Messenger,
) *UseCase {
	return &UseCase{
		logger:     logger,
		repository: repository,
		crawler:    crawler,
		messenger:  messenger,
	}
}
