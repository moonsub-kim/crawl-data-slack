package interpark

import (
	"fmt"
	"strconv"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(res Response, date string, crawlerName string, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event

	l := []string{"10:00", "10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30", "14:00", "14:30", "15:00", "15:30", "16:00", "16:30", "17:00"}
	startValue, err := strconv.Atoi(res.Data.RemainSeat[0].PlaySeq[1:])
	if err != nil {
		return nil, err
	}

	for _, seat := range res.Data.RemainSeat {
		if seat.RemainCnt == 0 {
			continue
		}

		t := seat.PlaySeq
		indexValue, err := strconv.Atoi(seat.PlaySeq[1:])
		if err == nil && indexValue-startValue < len(l) {
			t = l[indexValue-startValue]
		}

		msg := fmt.Sprintf("`%s %s` 생긴 티켓 %d", date, t, seat.RemainCnt)
		id := msg + time.Now().String()
		events = append(
			events,
			crawler.Event{
				Crawler:   crawlerName,
				Job:       jobName,
				UserName:  channel,
				UID:       id,
				Name:      id,
				EventTime: time.Now(), // TODO exact event time
				Message:   msg,
			},
		)
	}

	return events, nil
}
