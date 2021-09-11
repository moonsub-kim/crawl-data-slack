package hackernews

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
}

func (c Crawler) GetCrawlerName() string { return "hacker-news" }
func (c Crawler) GetJobName() string     { return "news" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var jsonBody string
	var dtos []DTO

	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate("https://news.ycombinator.com/news"),

		chromedp.Sleep(time.Second*2),

		// 문서 파싱
		chromedp.Evaluate(
			`
			function map_object(athing, metadata) {
				var a = athing.getElementsByClassName("storylink")[0];
				var subtext = metadata.getElementsByClassName("subtext")[0];

				// const regexHide = /(\| )?hide( \| )?/i;
				// var txt = subtext.innerText.replace(regex, '');
				// const regexAuthor = /(by .+ )/;
				// var txt = txt.replace(regex, '');

				return {
					"id": athing.id,
					"url": a.href,
					"comment_url": "https://news.ycombinator.com/item?id=" + athing.id,
					"title": a.innerText,
					"subtext": subtext.innerText,
				}
			}

			function crawl() {
				var trs = document.body.querySelectorAll('.itemlist > tbody > tr');
				var records = [];
				for (var i = 0; i < trs.length; i+=3) {
					var athing = trs[i];
					var metadata = trs[i+1];

					if (athing.classList.contains('athing') === false) {
						break;
					}

					records.push(map_object(athing, metadata));
				}

				return JSON.stringify(records);
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

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName())
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

func NewCrawler(logger *zap.Logger, chromectx context.Context) *Crawler {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,
	}
}
