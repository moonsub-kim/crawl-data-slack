package hackernews

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		// 292 points by geox 16 hours ago | hide | 140 comments
		// 49 minutes ago | hide
		// 11 points by todsacerdoti 46 minutes ago | hide | discuss

		age := regexp.MustCompile(`\d+ [A-z]+ ago`).FindString(dto.SubText)
		if strings.Contains(age, "minute") {
			fmt.Printf("ignore recent 1h post %v\n", dto)
			continue
		} else if !strings.Contains(dto.SubText, "comment") && !strings.Contains(dto.SubText, "discuss") {
			fmt.Printf("ignore ad post %v\n", dto)
			continue
		}

		comments := regexp.MustCompile(`\d+.comments?`).FindString(dto.SubText)
		comments = strings.ReplaceAll(comments, "&nbsp;", " ")
		if comments == "" {
			comments = "discuss"
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
					"<%s|%s>\n(%s, <%s|%s>)",
					dto.URL,
					dto.Title,
					age,
					dto.CommentURL,
					comments,
				),
			},
		)
	}

	return events, nil
}
