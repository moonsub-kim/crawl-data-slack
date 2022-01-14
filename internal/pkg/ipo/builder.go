package ipo

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(ipoList IPOMonthlyList, crawlerName, jobName string, channel string) []crawler.Event {
	events := b.parse(ipoList.LastMonthCompanies, crawlerName, jobName, channel)
	events = append(events, b.parse(ipoList.CurrentMonthCompanies, crawlerName, jobName, channel)...)
	events = append(events, b.parse(ipoList.NextMonthCompanies, crawlerName, jobName, channel)...)
	return events
}

func (b eventBuilder) parse(companies []Company, crawlerName, jobName string, channel string) []crawler.Event {
	var events []crawler.Event
	for _, d := range companies {
		if d.State != "공모청약" {
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      d.Name,
				Name:     d.Name,
				Message: fmt.Sprintf(
					"%s(%s) 상장예정! 공모청약기간 (%s ~ %s) <%s|dart>",
					d.Name,
					d.Code,
					d.StartDate,
					d.EndDate,
					d.Dart,
				),
			},
		)
	}

	return events
}
