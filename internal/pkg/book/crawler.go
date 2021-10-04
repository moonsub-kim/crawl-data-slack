package book

import (
	"net/url"
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

func (c Crawler) GetCrawlerName() string { return "kyobo" }
func (c Crawler) GetJobName() string     { return "apple" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(`https://search.kyobobook.co.kr/web/search?vPstrKeyWord=%25EC%2595%25A0%25ED%2594%258C&searchPcondition=1&searchPubNm=%EC%95%A0%ED%94%8C&searchCategory=%EA%B8%B0%ED%94%84%ED%8A%B8@GIFT&collName=GIFT&from_CollName=%EA%B8%B0%ED%94%84%ED%8A%B8@GIFT&searchOrder=5&vPstrTab=PRODUCT&from_coll=GIFT&currentPage=1&orderClick=LIe`)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	trs := doc.Find("tbody", "id", "search_list").FindAll("tr")
	var dtos []DTO
	for _, tr := range trs {
		td_detail := tr.Find("td", "class", "detail")
		a := td_detail.Find("div", "class", "title").Find("a")
		name := strings.TrimSpace(a.Find("strong").Text())
		uri := a.Attrs()["href"]
		url, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		id := url.Query().Get("barcode")
		price := tr.Find("td", "class", "price").Find("div", "class", "sell_price").Find("strong").Text()

		dtos = append(dtos, DTO{
			ID:    strings.TrimSpace(id),
			Name:  strings.TrimSpace(name),
			URL:   strings.TrimSpace(uri),
			Price: strings.TrimSpace(price),
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
