package quastor

import (
	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://blog.quastor.org/"
const BASE_URL string = "https://blog.quastor.org"
const iso8601Format = "2006-01-02T15:04:05.000Z"

func (c Crawler) GetCrawlerName() string { return "quastor" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	as := doc.Find("main", "class", "w-screen").
		// Find("div", "class", "max-w-none").
		// FindAll("div", "class", "px-4")[1].
		FindAll("a")

	var dtos []DTO
	for _, a := range as {
		if a.Find("svg").Error == nil { // premium contents 인경우
			continue
		}

		name := a.Find("h2").Text()
		url := BASE_URL + a.Attrs()["href"]
		date := a.Find("time").Attrs()["datetime"]

		dtos = append(dtos, DTO{
			Name: name,
			URL:  url,
			Date: date,
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
