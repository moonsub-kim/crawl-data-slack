package guardian

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

// month (jan, feb, mar, ...), day(1, 2, ...)
// const URL_TEMPLATE string = "https://www.theguardian.com/world/live/2022/%s/%d/russia-ukraine-latest-news-missile-strikes-on-oil-facilities-reported-as-some-russian-banks-cut-off-from-swift-system-live"
// https://www.theguardian.com/world/live/2022/feb/27/russia-ukraine-latest-news-missile-strikes-on-oil-facilities-reported-as-some-russian-banks-cut-off-from-swift-system-live

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel      string
	url_template string
}

func (c Crawler) GetCrawlerName() string { return "guardian" }
func (c Crawler) GetJobName() string     { return "guardian" }

func (c Crawler) Crawl() ([]crawler.Event, error) {

	now := time.Now()
	month := strings.ToLower(now.Month().String())[:3]
	day := now.Day()
	url := fmt.Sprintf(c.url_template, month, day)

	c.logger.Info(
		"URL",
		zap.Any("url", url),
	)

	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	divs := doc.FindAll("div", "class", "block--content")
	var dtos []DTO
	for _, div := range divs {
		c.logger.Info(
			"div",
			zap.Any("text", div.FullText()),
		)

		title := " "
		if h2 := div.Find("h2"); h2.Error == nil {
			title = h2.FullText()
		}

		content := div.Find("div", "class", "block-elements").FullText()
		re := regexp.MustCompile("( *\n)+")
		content = re.ReplaceAllString(content, "\n")

		dtos = append(
			[]DTO{ // reversed order
				{
					ID:       div.Attrs()["id"],
					Datetime: div.Find("time").Attrs()["datetime"],
					Title:    title,
					Content:  content,
				},
			},
			dtos...)
	}

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
	if err != nil {
		return nil, err
	}

	c.logger.Info(
		"crawler",
		zap.Any("events", events),
	)

	return events, nil
}

func NewCrawler(logger *zap.Logger, channel string, url_template string) *Crawler {
	return &Crawler{
		logger:       logger,
		channel:      channel,
		url_template: url_template,
	}
}
