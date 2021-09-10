package hackernews

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: "hacker-news",
				UID:      dto.ID,
				Name:     "news",
				Message:  fmt.Sprintf("<%s|%s> (<%s|comments>)", dto.URL, dto.Title, dto.CommentURL),
			},
		)
	}

	return events, nil
}
