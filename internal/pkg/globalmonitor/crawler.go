package globalmonitor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

// curl 'https://globalmonitor.einfomax.co.kr/bizrpt/reportlist' \
//   -H 'authority: globalmonitor.einfomax.co.kr' \
//   -H 'accept: application/json, text/plain, */*' \
//   -H 'accept-language: ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7' \
//   -H 'content-type: application/json;charset=UTF-8' \
//   -H 'origin: https://globalmonitor.einfomax.co.kr' \
//   -H 'referer: https://globalmonitor.einfomax.co.kr/ht_usa.html' \
//   -H 'sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="100", "Google Chrome";v="100"' \
//   -H 'sec-ch-ua-mobile: ?0' \
//   -H 'sec-ch-ua-platform: "macOS"' \
//   -H 'sec-fetch-dest: empty' \
//   -H 'sec-fetch-mode: cors' \
//   -H 'sec-fetch-site: same-origin' \
//   -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36' \
//   --data-raw '{"targetPeriodCheck":false,"page":1,"sortType":"latest","lscCd":700,"targetPeriodDate":null,"sscCd":"524910,521090,523050,523070","authSscCd":0,"cmd":"rl_011","startDate":"2022-01-20","searchItem":"title","endDate":"2022-04-20","searchStr":"","authLscCd":0,"pagePerItem":10}' \
//   --compressed

// const URL string = "https://globalmonitor.einfomax.co.kr/ht_usa.html#/3/01"
const URL string = "https://globalmonitor.einfomax.co.kr/bizrpt/reportlist"

func (c Crawler) GetCrawlerName() string { return "financial-report" }
func (c Crawler) GetJobName() string     { return "einfomax" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	reqBody := NewRequestBody(time.Now().AddDate(0, 0, -3), time.Now().AddDate(0, 0, 1))
	c.logger.Info("request_body", zap.Any("body", reqBody))
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("authority", "globalmonitor.einfomax.co.kr")
	header.Add("accept", "application/json, text/plain, */*")
	header.Add("accept-language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Add("content-type", "application/json;charset=UTF-8")
	header.Add("origin", "https://globalmonitor.einfomax.co.kr")
	header.Add("referer", "https://globalmonitor.einfomax.co.kr/ht_usa.html")
	header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
	header.Add("sec-ch-ua-mobile", "?0")
	header.Add("sec-ch-ua-platform", "\"macOS\"")
	header.Add("sec-fetch-dest", "empty")
	header.Add("sec-fetch-mode", "cors")
	header.Add("sec-fetch-site", "same-origin")
	header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36")

	buf := bytes.NewBuffer(body)

	res, err := http.Post(URL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
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
