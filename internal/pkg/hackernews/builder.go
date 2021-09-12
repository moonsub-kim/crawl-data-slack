package hackernews

import (
	"fmt"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) filter(filters []Filter, subText string) (string, string, bool) {
	var parsed string
	for _, f := range filters {
		if f.Filter(subText) {
			return f.Reason(), "", true
		}
		parsed += f.String() + " "
	}

	return "", strings.TrimSpace(parsed), false
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	filters := []Filter{&adFilter{}, &ageFilter{}, &pointFilter{}, &commentFilter{}}

	for _, dto := range dtos {
		// 292 points by geox 16 hours ago | hide | 140 comments
		// 49 minutes ago | hide
		// 11 points by todsacerdoti 46 minutes ago | hide | discuss
		reason, subText, filtered := b.filter(filters, dto.SubText)
		if filtered {
			fmt.Printf("%s %v\n", reason, dto)
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      dto.ID,
				Name:     "news",
				Message: fmt.Sprintf(
					"<%s|%s>\n(%s|%s>)",
					dto.URL,
					dto.Title,
					dto.CommentURL,
					subText,
				),
			},
		)
	}

	return events, nil
}
