package rss

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

const iso8601Format = "2006-01-02T15:04:05Z07:00"

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(feed *gofeed.Feed, crawlerName string, jobName string, transformers []Transformer, channel string) ([]crawler.Event, error) {
	var events []crawler.Event

	optionalTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return fmt.Sprintf("[%v]", t.Format(iso8601Format))
	}

	for i := len(feed.Items) - 1; i >= 0; i-- {
		item := b.transform(transformers, feed.Items[i])
		if item == nil {
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
					item.Description,
				),
			},
		)
	}

	return events, nil
}

func (b eventBuilder) transform(transformers []Transformer, item *gofeed.Item) *gofeed.Item {
	for _, t := range transformers {
		item = t.Transform(item)
		if item == nil {
			return nil
		}
	}

	return item
}
