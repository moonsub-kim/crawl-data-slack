package gitpublic

import (
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
}

func (c Crawler) GetCrawlerName() string { return "git" }
func (c Crawler) GetJobName() string     { return "public_repo" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(`https://github.com/orgs/Buzzvil/repositories`)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	lis := doc.Find("div", "class", "Box").Find("ul").FindAll("li")
	var dtos []DTO
	for _, li := range lis {
		a := li.Find("div", "class", "public").
			Find("div", "class", "flex-justify-between").
			Find("div", "class", "flex-auto").
			Find("a", "class", "d-inline-block")
		name := a.Text()
		url := "https://github.com" + a.Attrs()["href"]

		dtos = append(dtos, DTO{
			ID:   strings.TrimSpace(name),
			Name: strings.TrimSpace(name),
			URL:  strings.TrimSpace(url),
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
