package eomisae

import (
	"context"
	"encoding/json"
	"fmt"
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
	target  target
	id      string
	pw      string
}

const LOGIN_URL string = "https://eomisae.co.kr/index.php?act=dispMemberLoginForm"

func (c Crawler) GetCrawlerName() string { return "eomisae" }
func (c Crawler) GetJobName() string     { return c.target.name }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var body string
	var links []string
	var dtos []DTO

	err := chromedp.Run(
		c.ctx,

		// 로그인페이지: 로그인
		chromedp.Navigate(LOGIN_URL),
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`
				document.getElementById('uid').value = '%s';
				document.getElementById('upw').value = '%s';
				document.getElementsByClassName('submit')[0].click();
				`,
				c.id,
				c.pw,
			),
			nil,
		),
		chromedp.Sleep(time.Second*2),

		// go to url
		chromedp.Navigate(c.target.url),
		chromedp.EvaluateAsDevTools(
			`
			l = document.querySelectorAll('div.card_el > a.pjax');
			var links = [];
			for (var i = 0; i < l.length; i++) {
				links.push(l[i].href)
			}
			return JSON.stringify(records)
			`,
			&body,
		),
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &links)
	if err != nil {
		return nil, err
	}

	for _, l := range links {
		var dto DTO
		err := chromedp.Run(
			c.ctx,

			// link page
			chromedp.Navigate(l),
			chromedp.Sleep(time.Second*2),
			chromedp.EvaluateAsDevTools(
				c.target.script,
				&dto,
			),
		)
		if err != nil {
			c.logger.Error(
				"failed to run crawler",
				zap.Error(err),
				zap.Any("link", l),
				zap.Any("target", c.target.name),
			)
			continue
		}
		dtos = append(dtos, dto)
	}

	if len(dtos) == 1 {
		return []crawler.Event{}, nil
	}

	c.logger.Info("dto", zap.Any("dto", dtos))
	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewCrawler(logger *zap.Logger, chromectx context.Context, channel string, target string, id string, pw string) (*Crawler, error) {
	t, ok := targets[target]
	if !ok {
		return nil, fmt.Errorf("target %s not matched", target)
	}

	return &Crawler{
		logger: logger,
		ctx:    chromectx,

		channel: channel,
		target:  t,
		id:      id,
		pw:      pw,
	}, nil
}
