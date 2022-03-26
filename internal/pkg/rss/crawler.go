package rss

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	siteName string
	channel  string
	name     string

	filters []Filter
}

func (c Crawler) GetCrawlerName() string { return "rss" }
func (c Crawler) GetJobName() string     { return c.name }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	c.logger.Info(
		"site", zap.Any("site", c.siteName),
	)
	p := gofeed.NewParser()
	feed, err := p.ParseURLWithContext(c.siteName, ctx)
	if err != nil {
		return nil, err
	}

	return c.eventBuilder.buildEvents(feed, c.GetCrawlerName(), c.GetJobName(), c.filters, c.channel)
}

func NewCrawler(logger *zap.Logger, channel string, name string, siteName string, options ...CrawlerOption) *Crawler {
	c := &Crawler{
		logger:   logger,
		name:     name,
		siteName: siteName,
		channel:  channel,
	}

	for _, o := range options {
		o(c)
	}

	return c
}
