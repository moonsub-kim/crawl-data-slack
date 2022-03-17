package confluent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://docs.confluent.io/cloud/current/release-notes/index.html"

func (c Crawler) GetCrawlerName() string { return "confluent" }
func (c Crawler) GetJobName() string     { return "release" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var jsonBody string
	var dtos []DTO

	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate(URL),
		chromedp.Sleep(time.Second*2),
		chromedp.EvaluateAsDevTools(
			`
			function map_object(div) {
				return {
					"date": div.querySelector('h2').innerText,
					"content": div.innerText,
				};
			}

			function crawl() {
				var records = [];
				var divs = document.querySelectorAll('div#ccloud-release-notes > div.section')
				for (var i = 0; i < Math.min(divs.length, 5); i++) { // 최근 5개만 확인
					records.push(map_object(divs[i]));
				}
				
				return JSON.stringify(records.reverse()); // 과거 릴리즈가 인덱스 앞쪽으로 하기위해 reverse
			}
			crawl();
			`,
			&jsonBody,
		),
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonBody), &dtos)
	if err != nil {
		return nil, err
	}

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel, URL)
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

func NewCrawler(logger *zap.Logger, chromectx context.Context, channel string) *Crawler {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,

		channel: channel,
	}
}
