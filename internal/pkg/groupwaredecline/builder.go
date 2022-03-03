package groupwaredecline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dtos []DTO, crawlerName, jobName string, masters []string) ([]crawler.Event, error) {
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
				Message:  fmt.Sprintf("결재(`%s`)가 반려되었습니다. <https://gr.buzzvil.com/gw/userMain.do|그룹웨어>에서 확인해주세요.\n반려된 문서는 수정이 불가능하므로 새로 작성해주셔야 하며, 기존에 작성해둔 결제건에 대한 적요, 증빙유형등은 그대로 남아있습니다.", dto.DocName),
			},
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: reviewerName,
				UID:      dto.UID,
				Name:     "notified_declined",
				Message:  fmt.Sprintf("%s 에게 결재(`%s`) 반려 알림이 전달되었습니다.", drafterName, dto.DocName),
			},
		)

		for _, master := range masters {
			if reviewerName == master {
				continue
			}

			events = append(events, crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: master,
				UID:      dto.UID,
				Name:     fmt.Sprintf("notified_declined_master_%s", master),
				Message:  fmt.Sprintf("%s 가 %s 에게 `%s`를 반려 처리 하였습니다.", reviewerName, drafterName, dto.DocName),
			})
		}
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
	re := regexp.MustCompile(`(.+)\(.+/([A-z]+)[\. ]([A-z]+)\)`) // 상태(한글이름/firstname lastname) or 상태(한글이름/firstname.lastname)
	groups := re.FindStringSubmatch(status)
	if len(groups) != 4 {
		return "", "", fmt.Errorf("failed to parse status %s", status)
	}

	return groups[1], fmt.Sprintf("%s.%s", strings.ToLower(groups[2]), strings.ToLower(groups[3])), nil
}
