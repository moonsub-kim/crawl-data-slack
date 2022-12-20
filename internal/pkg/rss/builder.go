package rss

import (
	"fmt"
	"strings"
	"time"

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

	optionalTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return fmt.Sprintf("[%v] ", t.Format(iso8601Format))
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
				UID:      item.Link,
				Name:     item.Link,
				Message: fmt.Sprintf(
					"%s %s <%s|%s>\n%s",
					optionalTime(item.PublishedParsed),
					jobName,
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
