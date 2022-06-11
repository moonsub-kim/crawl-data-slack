package shopifyengineering

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

const URL string = "https://shopify.engineering"

func (c Crawler) GetCrawlerName() string { return "shopify-engineering" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) parseArticle(article soup.Root) DTO {
	a := article.Find("h2").Find("a")
	return DTO{
		ID:   a.Text(),
		Name: a.Text(),
		URL:  URL + a.Attrs()["href"],
		Date: article.Find("time").Attrs()["datetime"],
	}
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)

	divRoot := doc.Find("main", "id", "Main").
		Find("section", "class", "section").
		Find("div", "class", "grid").
		Find("div", "class", "grid__item")

	dtos := []DTO{c.parseArticle(divRoot.Find("article"))}

	for _, article := range divRoot.Find("div", "class", "grid--equal-height").FindAll("article") {
		dto := c.parseArticle(article)
		if !strings.HasPrefix(dto.Date, "2022-06-") {
			continue
		}
		dtos = append(dtos, dto)
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
