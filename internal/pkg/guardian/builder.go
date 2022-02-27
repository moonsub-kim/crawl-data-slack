package guardian

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		t := dto.UpdatedAt
		if dto.UpdatedAt != dto.CreatedAt {
			t = "Updated " + dto.UpdatedAt
		}
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      dto.ID + dto.UpdatedAt,
				Name:     dto.Title,
				Message: fmt.Sprintf(
					"[%s] *%s*\n%s",
					t,
					dto.Title,
					dto.Content,
				),
			},
		)
	}

	return events, nil
}
