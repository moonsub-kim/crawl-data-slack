package techcrunch

import (
	"fmt"

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
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      dto.CreatedAt,
				Name:     dto.Name,
				Message:  fmt.Sprintf("[%s] <%s|%s>", dto.CreatedAt, dto.URL, dto.Name),
			},
		)
	}

	return events, nil
}
