package eomisae

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
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

const LOGOUT_URL string = "https://eomisae.co.kr/index.php?mid=HO&act=dispMemberLogout"
const LOGIN_URL string = "https://eomisae.co.kr/index.php?act=dispMemberLoginForm"

func (c Crawler) GetCrawlerName() string { return "eomisae" }
func (c Crawler) GetJobName() string     { return c.target.name }

func (c Crawler) getLinks() ([]string, error) {
	res, err := soup.Get(c.target.url)
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	var links []string
	contents := doc.Find("div", "class", "bd_card").FindAll("a", "class", "hx")
	for _, content := range contents {
		links = append(links, content.Attrs()["href"])
	}

	return links, nil
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO

	links, err := c.getLinks()
	if err != nil {
		c.logger.Error(
			"failed to get links",
			zap.Error(err),
		)
		return nil, err
	} else if len(links) == 0 {
		c.logger.Info("no links to parse")
		return []crawler.Event{}, nil
	}
	c.logger.Info(
		"links",
		zap.Any("links", links),
	)

	// create actions
	bodies := make([]string, len(links))
	linkActions := []chromedp.Action{chromedp.Emulate(device.IPhone11)}
	for i, l := range links {
		linkActions = append(
			linkActions,
			c.createLinkActions(l, &bodies[i])...,
		)
	}

	actions := []chromedp.Action{
		chromedp.Emulate(device.IPhone11), // mobile emulation

		// 로그인페이지: 로그인
		chromedp.Navigate(LOGIN_URL),
		chromedp.Sleep(time.Second * 1),
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`document.getElementById('uid').value = '%s';
					document.getElementById('upw').value = '%s';
					document.getElementsByClassName('submit')[0].click();`,
				c.id,
				c.pw,
			),
			nil,
		),
		chromedp.Sleep(time.Second * 1),
	}
	actions = append(actions, linkActions...)

	err = chromedp.Run(
		c.ctx,
		actions...,
	)
	if err != nil {
		c.logger.Error("run error", zap.Error(err))
		return nil, err
	}

	// unmarshalling
	for _, body := range bodies {
		var dto DTO
		err = json.Unmarshal([]byte(body), &dto)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}

	if len(dtos) == 1 {
		return []crawler.Event{}, nil
	}

	c.logger.Info("dto", zap.Any("dto", dtos))
	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func (c Crawler) createLinkActions(link string, body *string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.ActionFunc(func(context.Context) error {
			c.logger.Info(
				"run",
				zap.Any("link", link),
			)
			return nil
		}),
		chromedp.Navigate(link),
		chromedp.Sleep(time.Second * 1),
		chromedp.EvaluateAsDevTools(
			c.target.script,
			body,
		),
	}
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
