package navercareer

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) []crawler.Event {
	var events []crawler.Event
	for _, d := range dtos {
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      d.Title,
				Name:     d.Title,
				Message: fmt.Sprintf(
					"<%s|%s>\n> %s",
					d.URL,
					d.Title,
					d.Info,
				),
			},
		)
	}

	return events
}
