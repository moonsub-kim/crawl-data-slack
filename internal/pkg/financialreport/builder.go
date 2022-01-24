package financialreport

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
				UID:      dto.Title,
				Name:     jobName,
				Message:  fmt.Sprintf("<%s|%s> (%s)\n%s", dto.Title, dto.PDFURL, dto.Company, dto.Text),
			},
		)
	}

	return events, nil
}
