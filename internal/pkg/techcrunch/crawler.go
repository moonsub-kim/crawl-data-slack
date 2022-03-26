package techcrunch

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

const URL string = "https://techcrunch.com/"

func (c Crawler) GetCrawlerName() string { return "techcrunch" }
func (c Crawler) GetJobName() string     { return "news" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	divs := doc.Find("div", "class", "river--homepage").FindAll("div", "class", "post-block")
	var dtos []DTO
	for _, div := range divs {
		a := div.Find("a", "class", "post-block__title__link")
		name := a.Text()
		url := a.Attrs()["href"]
		time := div.Find("time", "class", "river-byline__time")

		dtos = append(dtos, DTO{
			ID:        strings.TrimSpace(name),
			Name:      strings.TrimSpace(name),
			URL:       strings.TrimSpace(url),
			CreatedAt: time.Attrs()["datetime"],
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
