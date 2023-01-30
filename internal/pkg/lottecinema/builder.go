package lottecinema

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(items []Item, crawlerName, jobName string, channel string) []crawler.Event {
	var events []crawler.Event
	for _, i := range items {
		id := fmt.Sprintf("%s %s %s", i.PlayDt, i.StartTime, i.ScreenNameKR)
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      id,
				Name:     id,
				Message: fmt.Sprintf(
					"%s %s %s // 자리 %d/%d",
					i.FilmNameKR,
					i.ScreenDivisionNameKR,
					id,
					i.BookingSeatCount,
					i.TotalSeatCount,
				),
			},
		)
	}

	return events
}
