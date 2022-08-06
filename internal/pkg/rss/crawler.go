package rss

import (
	"io"
	"net/http"
	"regexp"
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
	c.logger.Info(
		"site", zap.Any("site", c.siteName),
	)

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(c.siteName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(string(rune(8))) // remove backspace charater
	str := pattern.ReplaceAllString(string(b), "")

	parser := gofeed.NewParser()
	feed, err := parser.ParseString(str)
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
