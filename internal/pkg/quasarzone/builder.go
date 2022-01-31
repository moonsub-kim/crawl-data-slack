package quasarzone

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) filter(filters []Filter, dto DTO) (string, bool) {
	for _, f := range filters {
		if f.Filter(dto) {
			return f.Reason(), true
		}
	}

	return "", false
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	filters := []Filter{&statusFilter{}}

	for _, dto := range dtos {
		reason, filtered := b.filter(filters, dto)
		if filtered {
			fmt.Printf("%s %v\n", reason, dto)
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      dto.ID,
				Name:     "sales-pc",
				Message: fmt.Sprintf(
					"[%s] <%s|%s>\n%s (%s)",
					dto.Status,
					dto.URL,
					dto.Name,
					dto.PriceInfo,
					dto.Date,
				),
			},
		)
	}

	return events, nil
}
