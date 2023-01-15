package designerjob

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(res Response, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event

	for _, d := range res.Data {
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      fmt.Sprintf("%s-%s", d.Company.Name, d.Position),
				Name:     "position",
				Message: fmt.Sprintf(
					"[%s] %s\n(%s)",
					d.Company.Name,
					d.Position,
					fmt.Sprintf("https://www.wanted.co.kr/wd/%d", d.ID),
				),
			},
		)
	}

	return events, nil
}
