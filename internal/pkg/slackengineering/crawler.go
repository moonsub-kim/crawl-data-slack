package slackengineering

import (
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://slack.engineering"

func (c Crawler) GetCrawlerName() string { return "slack-engineering" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	divs := doc.Find("div", "class", "post-list").FindAll("h3", "class", "post-item__title")
	var dtos []DTO
	for _, div := range divs {
		a := div.Find("a")
		name := a.Text()
		url := a.Attrs()["href"]

		dtos = append(dtos, DTO{
			ID:   strings.TrimSpace(name),
			Name: strings.TrimSpace(name),
			URL:  strings.TrimSpace(url),
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
