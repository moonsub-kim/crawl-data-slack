package wanted

import (
	"fmt"
	"strings"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(res Response, crawlerName, jobName string, channel string, excludes []string) ([]crawler.Event, error) {
	var events []crawler.Event

	for _, d := range res.Data {
		if b.isExcludes(excludes, strings.ToLower(d.Company.Name)) {
			continue
		}

		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       fmt.Sprintf("%s-%s", d.Company.Name, d.Position),
				Name:      "position",
				EventTime: time.Now(), // TODO exact event time
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

func (b eventBuilder) isExcludes(excludes []string, company string) bool {
	for _, s := range excludes {
		if strings.Contains(company, s) {
			return true
		}
	}
	return false
}
