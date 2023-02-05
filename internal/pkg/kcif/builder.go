package kcif

import (
	"fmt"
	"time"

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
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       d.ID,
				Name:      d.Title,
				EventTime: time.Now(), // TODO exact event time
				Message: fmt.Sprintf(
					"[%s] 국제금융센터 *%s*, <%s|PDF 보기>\n> %s",
					d.Date,
					d.Title,
					d.pdfURL,
					d.Content,
				),
			},
		)
	}

	return events
}
