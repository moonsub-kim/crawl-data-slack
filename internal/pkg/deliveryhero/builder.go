package deliveryhero

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

const DATE_PARSE_FORMAT string = "02.01.06"

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		name := dto.Name
		if len(name) > 64 {
			name = name[:64]
		}

		date, err := time.Parse(DATE_PARSE_FORMAT, dto.Date)
		if err != nil {
			return nil, err
		}

		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       name,
				Name:      name,
				EventTime: date,
				Message:   fmt.Sprintf("DeliveryHero <%s|%s>", dto.URL, dto.Name),
			},
		)
	}

	return events, nil
}
