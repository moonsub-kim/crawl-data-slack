package slackengineering

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		name := dto.Name
		if len(name) > 64 {
			name = name[:64]
		}
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      name,
				Name:     name,
				Message:  fmt.Sprintf("<%s|%s>", dto.URL, dto.Name),
			},
		)
	}

	return events, nil
}
