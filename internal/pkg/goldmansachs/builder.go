package goldmansachs

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

const DATE_PARSE_FORMAT string = "January 2 2006"

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event

	for _, dto := range dtos {
		eventTime, err := time.Parse(DATE_PARSE_FORMAT, dto.Date)
		if err != nil {
			return nil, err
		}

		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       dto.Name,
				Name:      dto.Name,
				EventTime: eventTime,
				Message:   fmt.Sprintf("[%s] GoldmanSachs <%s|%s>", dto.Date, dto.URL, dto.Name),
			},
		)
	}

	return events, nil
}
