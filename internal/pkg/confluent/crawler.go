package confluent

import (
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel      string
	organization string
}

const URL string = "https://docs.confluent.io/cloud/current/release-notes/index.html"

func (c Crawler) GetCrawlerName() string { return "confluent" }
func (c Crawler) GetJobName() string     { return "release" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	divs := doc.Find("div", "id", "ccloud-release-notes").FindAll("div", "class", "section")
	var dtos []DTO
	for _, div := range divs[:5] { // 최근5개만 확인
		date := div.Find("h2").Text()
		content := div.Text()

		dtos = append(dtos, DTO{
			Date:    strings.TrimSpace(date),
			Content: strings.TrimSpace(content),
		})
	}

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
	if err != nil {
		return nil, err
	}

	c.logger.Info(
		"crawler",
		// zap.Any("dto", dtos),
		zap.Any("events", events),
	)
	return events, nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
