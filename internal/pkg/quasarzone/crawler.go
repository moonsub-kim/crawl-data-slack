package quasarzone

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

var regexID *regexp.Regexp = regexp.MustCompile(`views/\d+`)
var regexPrice *regexp.Regexp = regexp.MustCompile(`^.+? 가격 +(.+)`)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
}

func (c Crawler) GetCrawlerName() string { return "quasarzone" }
func (c Crawler) GetJobName() string     { return "sale-pc" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO
	res, err := soup.Get("https://quasarzone.com/bbs/qb_saleinfo?_method=post&type=&page=1&category=PC%2F%ED%95%98%EB%93%9C%EC%9B%A8%EC%96%B4&popularity=&kind=subject&keyword=&sort=num%2C+reply&direction=DESC")
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	contents := doc.FindAll("div", "class", "market-info-list-cont")
	for _, content := range contents {
		dto := c.crawlEntry(content)
		if dto.isEmpty() {
			continue
		}

		dtos = append([]DTO{dto}, dtos...) // insert reversed order
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

func (c Crawler) crawlEntry(content soup.Root) DTO {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Panic", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	div_market_info_sub_ps := content.Find("div", "class", "market-info-sub").FindAll("p")
	p_tit := content.Find("p", "class", "tit")
	url := "https://quasarzone.com" + p_tit.Find("a", "class", "subject-link").Attrs()["href"]

	spans := div_market_info_sub_ps[0].FindAll("span")
	var s string
	for _, span := range spans {
		s += span.Text() + " "
	}

	return DTO{
		ID:        strings.TrimSpace(regexID.FindString(url)[6:]),
		Status:    strings.TrimSpace(p_tit.Find("span", "class", "label").Text()),
		Name:      strings.TrimSpace(p_tit.Find("a", "class", "subject-link").Find("span", "class", "ellipsis-with-reply-cnt").Text()),
		URL:       strings.TrimSpace(url),
		Category:  strings.TrimSpace(div_market_info_sub_ps[0].Find("span", "class", "category").Text()),
		PriceInfo: strings.TrimSpace(regexPrice.FindStringSubmatch(s)[1]),
		Date:      strings.TrimSpace(div_market_info_sub_ps[1].Find("span", "class", "date").Text()),
	}
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
