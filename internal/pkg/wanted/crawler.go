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
	excludes     []string
}

func (c Crawler) GetCrawlerName() string { return "wanted" }
func (c Crawler) GetJobName() string     { return "open-position" }

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

	return c.eventBuilder.buildEvents(response, c.GetCrawlerName(), c.GetJobName(), c.channel, c.excludes)
}

func NewCrawler(logger *zap.Logger, channel string, query string, excludes []string) *Crawler {
	return &Crawler{
		logger:   logger,
		channel:  channel,
		query:    url.QueryEscape(query),
		excludes: excludes,
	}
}
