package miraeasset

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
				UID:      d.ID,
				Name:     d.Title,
				Message:  b.buildMessage(d),
			},
		)
	}

	return events
}

func (b eventBuilder) buildMessage(d DTO) string {
	pdf := ""
	if d.pdfURL != nil {
		pdf = fmt.Sprintf(", <%s|PDF 보기>", *d.pdfURL)
	}
	return fmt.Sprintf(
		"*%s*\n%s %s%s\n> %s",
		d.Title,
		"미래에셋증권",
		d.Date,
		pdf,
		d.Content,
	)
}
