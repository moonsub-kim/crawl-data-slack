package datablog

import (
	"fmt"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type BlogCrawler interface {
	Crawl() ([]DTO, error)
}

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
}

func (c Crawler) GetCrawlerName() string { return "hacker-news" }
func (c Crawler) GetJobName() string     { return "news" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	return nil, fmt.Errorf("Unimplementd")
}

func NewCrawler(logger *zap.Logger, channel string, blogCrawler BlogCrawler) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
