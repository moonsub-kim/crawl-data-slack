package miraeasset

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
	query        string
}

const CONTENT_URL_TEMPLATE string = "https://securities.miraeasset.com/bbs/board/message/view.do?messageId=%s&messageNumber=%s&categoryId=1521"
const PDF_URL_TEMPLATE string = "https://securities.miraeasset.com/bbs/download/%s.pdf?attachmentId=%s"

func (c Crawler) GetCrawlerName() string { return "ipo" }
func (c Crawler) GetJobName() string     { return "ipo" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO

	res, err := soup.Get("https://securities.miraeasset.com/bbs/board/message/list.do?categoryId=1521")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	trs := doc.Find("table", "class", "bbs_linetype2").Find("tbody").FindAll("tr")
	for _, tr := range trs {
		tds := tr.FindAll("td")
		date := strings.TrimSpace(tds[0].Text())
		a := tds[1].Find("div", "class", "subject").Find("a")

		title := strings.ReplaceAll(strings.TrimSpace(a.FullText()), "\n", " ")
		parsed := strings.Split(a.Attrs()["href"], "'")
		id := parsed[1]
		number := parsed[3]
		contentURL := fmt.Sprintf(CONTENT_URL_TEMPLATE, id, number)
		pdfURL := ""
		if tds[2].Find("p").Error == nil {
			pdfURL = fmt.Sprintf(PDF_URL_TEMPLATE, id, id)
		} // pdfURL이 있는경우만

		content, err := c.parseContent(contentURL)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, DTO{
			ID:      id,
			Date:    date,
			Title:   title,
			pdfURL:  pdfURL,
			URL:     contentURL,
			Content: content,
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func (c Crawler) parseContent(contentURL string) (string, error) {
	res, err := soup.Get(contentURL)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(res)
	m := regexp.MustCompile(`\n+`)
	return m.ReplaceAllString(doc.Find("div", "id", "messageContentsDiv").FullText(), "\n"), nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
