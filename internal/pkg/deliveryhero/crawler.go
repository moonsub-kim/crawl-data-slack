package deliveryhero

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://tech.deliveryhero.com/"
const BASE_URL string = "https://tech.deliveryhero.com"

func (c Crawler) GetCrawlerName() string { return "delivery-hero" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var dtos []DTO
	doc.Find("div.container > article.blog-article > a").Each(
		func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			url := BASE_URL + href
			name := s.Find("h5").Text()
			date := s.Find("span.mr-1").Text()

			dtos = append(dtos, DTO{
				ID:   strings.TrimSpace(name),
				Name: strings.TrimSpace(name),
				URL:  strings.TrimSpace(url),
				Date: date,
			})
		},
	)
	c.logger.Info("dto", zap.Any("dtos", dtos))
	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
