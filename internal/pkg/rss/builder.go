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

	optionalString := func(s string) string {
		if s == "" {
			return ""
		}
		return fmt.Sprintf("(%s)", s)
	}

	for i := len(feed.Items) - 1; i >= 0; i-- {
		item := feed.Items[i]
		if b.filter(filters, item) {
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      item.Title,
				Name:     item.Title,
				Message: fmt.Sprintf(
					"[%v] <%s|%s>\n%s",
					item.PublishedParsed.Format(iso8601Format),
					item.Link,
					item.Title,
					optionalString(strings.Join(item.Categories, ", ")),
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
