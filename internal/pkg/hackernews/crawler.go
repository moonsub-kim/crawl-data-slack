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
			function map_object(tr) {
				var a = tr.getElementsByClassName("storylink")[0];
				return {
					"id": tr.id,
					"url": a.href,
					"comment_url": "https://news.ycombinator.com/item?id=" + tr.id,
					"title": a.innerText,
				}
			}

			function crawl() {
				var trs = document.body.querySelectorAll('.itemlist > tbody > tr.athing');
				var records = [];
				if (trs.length == 0) {
					return "[]" 			// ignore empty search results
				}

				for (var i = 0; i < trs.length; i++) {
					records.push(map_object(trs[i]));
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

	c.logger.Info("dto", zap.Any("dto", dtos))
	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName())
	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewCrawler(logger *zap.Logger, chromectx context.Context) *Crawler {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,
	}
}
