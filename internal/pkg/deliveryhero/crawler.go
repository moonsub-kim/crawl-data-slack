package deliveryhero

import (
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel    string
	recentDays int64
}

const URL string = "https://tech.deliveryhero.com/"
const BASE_URL string = "https://tech.deliveryhero.com"

func (c Crawler) GetCrawlerName() string { return "delivery-hero" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	as := doc.Find("div", "class", "container").FindAll("a")
	var dtos []DTO
	for _, a := range as {
		url := BASE_URL + a.Attrs()["href"]
		name := a.Find("h4").Text()

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
