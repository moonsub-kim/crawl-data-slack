package interpark

import (
	"fmt"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(res Response, crawlerName, jobName string, channel string) []crawler.Event {
	var events []crawler.Event
	for _, seat := range res.Data.RemainSeat {
		if seat.RemainCnt == 0 {
			continue
		}
		m := map[string]string{
			"I11": "10:00",
			"I12": "10:30",
			"I13": "11:00",
			"I14": "11:30",
			"I15": "12:00",
			"I16": "12:30",
			"I17": "13:00",
			"I18": "13:30",
			"I19": "14:00",
			"I20": "14:30",
			"I21": "15:00",
			"I22": "15:30",
			"I23": "16:00",
			"I24": "16:30",
			"I25": "17:00",
		}
		t, ok := m[seat.PlaySeq]
		if !ok {
			t = seat.PlaySeq
		}
		id := fmt.Sprintf("%s %s, 생긴 티켓 %d, %v", res.Common.RequestURI, t, seat.RemainCnt, time.Now())
		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       id,
				Name:      id,
				EventTime: time.Now(), // TODO exact event time
				Message:   id,
			},
		)
	}

	return events
}
