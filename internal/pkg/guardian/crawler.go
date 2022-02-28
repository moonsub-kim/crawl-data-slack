package guardian

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

// month (jan, feb, mar, ...)
// day(1, 2, ...)
// topic
const URL_TEMPLATE string = "https://www.theguardian.com/world/live/2022/%s/%d/%s"

// https://www.theguardian.com/world/live/2022/feb/27/russia-ukraine-latest-news-missile-strikes-on-oil-facilities-reported-as-some-russian-banks-cut-off-from-swift-system-live

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder

	channel string
	url     string
}

func (c Crawler) GetCrawlerName() string { return "guardian" }
func (c Crawler) GetJobName() string     { return "guardian" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	c.logger.Info(
		"URL",
		zap.Any("url", c.url),
	)

	var dtos []DTO
	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate(c.url),
		chromedp.Sleep(time.Second*3),
		chromedp.Evaluate(
			`
			function getTitle(div) {
				h2 = div.querySelector('h2');
				if (h2 === null) return ' ';
				return h2.innerText;
			}

			function getUpdatedAt(div) {
				p = div.querySelector('p.updated-time > time');
				if (p === null) return '';
				return p.getAttribute('datetime');
			}

			function main() {
				var dtos = [];
				divs = document.querySelectorAll('div.js-liveblog-body > div.is-key-event');
				for (let div of divs) {
					dtos.push({
						title: getTitle(div),
						id: div.id,
						created_at: div.querySelector('time').getAttribute('datetime'),
						updated_at: getUpdatedAt(div),
						content: div.querySelector('div.block-elements').innerText,
					});
				}
				return dtos;
			}

			main()
			`,
			&dtos,
		),
	)
	if err != nil {
		return nil, err
	}

	c.logger.Info(
		"crawler",
		zap.Any("dtos", dtos),
	)

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

func NewCrawler(logger *zap.Logger, chromectx context.Context, channel string, url string) *Crawler {
	return &Crawler{
		logger:  logger,
		ctx:     chromectx,
		channel: channel,
		url:     url,
	}
}
