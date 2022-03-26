package ipo

import (
	"encoding/json"

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

func (c Crawler) GetCrawlerName() string { return "ipo" }
func (c Crawler) GetJobName() string     { return "ipo" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var root Root

	res, err := soup.Get("https://www.ustockplus.com/ipo/calander#monthList")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	content := doc.Find("body").Find("script", "id", "__NEXT_DATA__")
	err = json.Unmarshal([]byte(content.Text()), &root)
	if err != nil {
		return nil, err
	}

	return c.eventBuilder.buildEvents(root.Props.PageProps.IPOMonthlyList, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
