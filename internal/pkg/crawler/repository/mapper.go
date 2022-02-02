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

func (m mapper) mapChannelToModelChannel(c crawler.Channel) Channel {
	return Channel{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (m mapper) mapChannelsToModelChannels(channels []crawler.Channel) []Channel {
	var modelChannels []Channel
	for _, u := range channels {
		modelChannels = append(modelChannels, m.mapChannelToModelChannel(u))
	}
	return modelChannels
}

func (m mapper) mapModelChannelToChannel(c Channel) crawler.Channel {
	return crawler.Channel{
		ID:   c.ID,
		Name: c.Name,
	}
}
