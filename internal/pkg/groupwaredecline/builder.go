package groupwaredecline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, dto := range dtos {
		drafterName, err := b.parseDrafter(dto.Drafter)
		if err != nil {
			return nil, err
		}

		_, reviewerName, err := b.parseStatus((dto.Status))
		if err != nil {
			return nil, err
		}

		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: drafterName,
				UID:      dto.UID,
				Name:     "declined",
				Message:  fmt.Sprintf("결재(%s)가 반려되었습니다. <https://gr.buzzvil.com/gw/userMain.do|그룹웨어>에서 확인해주세요.", dto.DocName),
			},
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: reviewerName,
				UID:      dto.UID,
				Name:     "notified_declined",
				Message:  fmt.Sprintf("%s 에게 결재(%s) 반려 알림이 전달되었습니다.", drafterName, dto.DocName),
			},
		)
	}

	return events, nil
}

func (b eventBuilder) parseDrafter(drafter string) (string, error) {
	re := regexp.MustCompile(`.+/([A-z]+) ([A-z]+)`) // 한글이름/firstname lastname
	groups := re.FindStringSubmatch(drafter)
	if len(groups) != 3 {
		return "", fmt.Errorf("failed to parse drafter %s", drafter)
	}

	return fmt.Sprintf("%s.%s", strings.ToLower(groups[1]), strings.ToLower(groups[2])), nil
}

func (b eventBuilder) parseStatus(status string) (string, string, error) {
	re := regexp.MustCompile(`(.+)\(.+/([A-z]+) ([A-z]+)\)`) // 상태(한글이름/firstname lastname)
	groups := re.FindStringSubmatch(status)
	if len(groups) != 4 {
		return "", "", fmt.Errorf("failed to parse status %s", status)
	}

	return groups[1], fmt.Sprintf("%s.%s", strings.ToLower(groups[2]), strings.ToLower(groups[3])), nil
}
