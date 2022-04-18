package financialreport

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder

	channel string
}

// const URL string = "https://globalmonitor.einfomax.co.kr/ht_usa.html#/3/01"
const URL string = "https://globalmonitor.einfomax.co.kr/bizrpt/reportlist"

func (c Crawler) GetCrawlerName() string { return "financial-report" }
func (c Crawler) GetJobName() string     { return "einfomax" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	ReqBody := NewRequestBody(time.Now().AddDate(0, 0, -3), time.Now().AddDate(0, 0, 1))
	body, err := json.Marshal(ReqBody)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(URL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var dto DTO
	err = json.Unmarshal(resBody, &dto)
	if err != nil {
		return nil, err
	}

	return c.eventBuilder.buildEvents(dto, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string) (*Crawler, error) {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}, nil
}
