package naverd2

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(contents []Content, crawlerName, jobName string, channel string, baseURL string) []crawler.Event {
	var events []crawler.Event
	for _, c := range contents {
		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       c.Path,
				Name:      c.Title,
				EventTime: time.Unix(c.PublishedAt/1000, 0),
				Message: fmt.Sprintf(
					"[Naver D2] <%s|%s>",
					baseURL+c.Path,
					c.Title,
				),
			},
		)
	}

	return events
}
