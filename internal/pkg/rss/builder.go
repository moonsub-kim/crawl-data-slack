package rss

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

const iso8601Format = "2006-01-02T15:04:05Z07:00"

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(feed *gofeed.Feed, crawlerName, jobName string, filters []Filter, channel string) ([]crawler.Event, error) {
	var events []crawler.Event

	for _, i := range feed.Items {
		if b.filter(filters, i) {
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      i.Title,
				Name:     i.Title,
				Message: fmt.Sprintf(
					"[%v] %s <%s|%s>\n(%s)",
					i.PublishedParsed.Format(iso8601Format),
					jobName,
					i.Link,
					i.Title,
					strings.Join(i.Categories, ", "),
				),
			},
		)
	}

	return events, nil
}

func (b eventBuilder) filter(filters []Filter, item *gofeed.Item) bool {
	for _, f := range filters {
		if f.Filter(item) {
			return true
		}
	}

	return false
}
