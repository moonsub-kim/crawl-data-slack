package quastor

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
				UID:       dto.Name,
				Name:      dto.Name,
				EventTime: time.Now(),
				Message:   fmt.Sprintf("[%s] Quastor <%s|%s>", dto.Date, dto.URL, dto.Name),
			},
		)
	}

	return events, nil
}
