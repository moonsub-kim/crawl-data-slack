package quasarzone

import (
	"regexp"
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

func (c Crawler) GetCrawlerName() string { return "quasarzone" }
func (c Crawler) GetJobName() string     { return "sale-pc" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO
	res, err := soup.Get("https://quasarzone.com/bbs/qb_saleinfo")
	if err != nil {
		return nil, err
	}

	reID := regexp.MustCompile(`\d+$`)
	rePrice := regexp.MustCompile(`^.+? 가격 +(.+)`)

	doc := soup.HTMLParse(res)
	contents := doc.FindAll("div", "class", "market-info-list-cont")
	for _, content := range contents {
		div_market_info_sub_ps := content.Find("div", "class", "market-info-sub").FindAll("p")
		p_tit := content.Find("p", "class", "tit")
		url := "https://quasarzone.com" + p_tit.Find("a", "class", "subject-link").Attrs()["href"]

		spans := div_market_info_sub_ps[0].FindAll("span")
		var s string
		for _, span := range spans {
			s += span.Text() + " "
		}
		dtos = append([]DTO{{
			ID:        strings.TrimSpace(reID.FindString(url)),
			Status:    strings.TrimSpace(p_tit.Find("span", "class", "label").Text()),
			Name:      strings.TrimSpace(p_tit.Find("a", "class", "subject-link").Find("span", "class", "ellipsis-with-reply-cnt").Text()),
			URL:       strings.TrimSpace(url),
			Category:  strings.TrimSpace(div_market_info_sub_ps[0].Find("span", "class", "category").Text()),
			PriceInfo: strings.TrimSpace(rePrice.FindStringSubmatch(s)[1]),
			Date:      strings.TrimSpace(div_market_info_sub_ps[1].Find("span", "class", "date").Text()),
		}}, dtos...) // insert reversed order
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
