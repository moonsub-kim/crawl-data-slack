package hackernews

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		// 292 points by geox 16 hours ago | hide | 140 comments
		// 49 minutes ago | hide
		// 11 points by todsacerdoti 46 minutes ago | hide | discuss
		m := getParams(`(?P<points>\d+)?( points?)?( by .+? )?(?P<age>\d+ .+ ago) \| hide( \| )?(?P<comments>\d+)?( comments?)?( discuss)?`, dto.SubText)

		// ignore recent 1h post
		if strings.Contains(m["age"], "minute") {
			fmt.Printf("ignore %v\n", dto)
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: "hacker-news",
				UID:      dto.ID,
				Name:     "news",
				Message:  fmt.Sprintf("<%s|%s>\n(%s <%s|%s comments>)\n", dto.URL, dto.Title, m["age"], dto.CommentURL, m["comments"]),
			},
		)
	}

	return events, nil
}

func getParams(regEx, s string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
