package confluent

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string, url string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       dto.Date,
				Name:      dto.Date,
				EventTime: time.Now(),
				Message:   fmt.Sprintf("%s\n<%s|RELEASE NOTE>", dto.Content, url),
			},
		)
	}

	return events, nil
}
