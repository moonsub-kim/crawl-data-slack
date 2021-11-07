package wanted

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
	query        string
}

func (c Crawler) GetCrawlerName() string { return "quasarzone" }
func (c Crawler) GetJobName() string     { return "sale-pc" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	url := fmt.Sprintf(
		"https://www.wanted.co.kr/api/v4/jobs?%d&country=kr&job_sort=company.response_rate_order&locations=all&years=-1&query=%s&limit=%d",
		time.Now().Unix(),
		c.query,
		100,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response Response
	json.NewDecoder(res.Body).Decode(&response)

	events, err := c.eventBuilder.buildEvents(response, c.GetCrawlerName(), c.GetJobName(), c.channel)
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

func NewCrawler(logger *zap.Logger, channel string, query string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
		query:   url.QueryEscape(query),
	}
}
