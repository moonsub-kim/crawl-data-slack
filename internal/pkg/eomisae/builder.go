package eomisae

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       dto.URL,
				Name:      jobName,
				EventTime: time.Now(), // TODO exact event time
				Message:   fmt.Sprintf("<%s|%s>\n%s\n%s", dto.Post, dto.Name, dto.URL, dto.Content),
			},
		)
	}

	return events, nil
}
