package financialreport

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
)

type eventBuilder struct {
}

func (b eventBuilder) buildEvents(dto DTO, crawlerName, jobName string, channel string) ([]crawler.Event, error) {
	var events []crawler.Event
	for _, report := range dto.ReportList {
		events = append(
			events,
			crawler.Event{
				Crawler:  crawlerName,
				Job:      jobName,
				UserName: channel,
				UID:      report.Title,
				Name:     jobName,
				Message:  b.buildMessage(report),
			},
		)
	}

	return events, nil
}

func (b eventBuilder) buildMessage(report Report) string {
	return fmt.Sprintf(
		"*%s*\n%s %s, <%s|PDF 보기>\n> %s",
		report.Title,
		report.Company,
		report.Date,
		fmt.Sprintf("https://rreport.einfomax.co.kr/report/%s.pdf", report.ID),
		report.Text,
	)
	// m := map[string]interface{}{
	// 	"blocks": []slack.Block{
	// 		slack.NewHeaderBlock(
	// 			slack.NewTextBlockObject(slack.PlainTextType, dto.Title, false, false),
	// 		),
	// 		slack.NewDividerBlock(),
	// 		slack.NewContextBlock(
	// 			"",
	// 			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("%s %s , <%s|PDF 보기>", dto.Company, dto.Date, dto.PDFURL), false, false),
	// 			slack.NewTextBlockObject(slack.MarkdownType, dto.Text, false, false),
	// 		),
	// 	},
	// }
	// bytes, err := json.Marshal(m)
	// if err != nil {
	// 	return "", err
	// }

	// return string(bytes), nil
}
