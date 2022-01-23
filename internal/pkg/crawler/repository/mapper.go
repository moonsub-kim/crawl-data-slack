package repository

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type mapper struct {
}

func (m mapper) mapEventToModelEvent(e crawler.Event) Event {
	return Event{
		Crawler:  e.Crawler,
		Job:      e.Job,
		UserName: e.UserName,
		UID:      e.UID,
		Name:     e.Name,
		Message:  e.Message,
	}
}

func (m mapper) mapChannelToModelChannel(u crawler.Channel) Channel {
	return Channel{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (m mapper) mapChannelsToModelChannels(users []crawler.Channel) []Channel {
	var modelUsers []Channel
	for _, u := range users {
		modelUsers = append(modelUsers, m.mapChannelToModelChannel(u))
	}
	return modelUsers
}

func (m mapper) mapModelChannelToChannel(u Channel) crawler.Channel {
	return crawler.Channel{
		ID:   u.ID,
		Name: u.Name,
	}
}
