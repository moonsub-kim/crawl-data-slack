package confluent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder

	channel  string
	jobName  string
	keywords []string
}

func (c Crawler) GetCrawlerName() string { return "confluent" }
func (c Crawler) GetJobName() string     { return c.jobName }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	switch c.GetJobName() {
	case "release":
		return c.CrawlRelease()
	case "status":
		return c.CrawlStatus()
	}
	return nil, fmt.Errorf("unsupported crawler %s", c.GetJobName())
}

func (c Crawler) CrawlStatus() ([]crawler.Event, error) {
	url := "https://status.confluent.cloud/"
	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	divs := doc.FindAll("div", "class", "unresolved-incident")
	if len(divs) == 0 { // empty incident
		return []crawler.Event{}, nil
	}

	pattern := regexp.MustCompile(`(\n\s+)+`)
	var events []crawler.Event
	for _, div := range divs {
		if c.containsKeyword(div.FullText()) {
			a := div.Find("a", "class", "actual-title")
			text := div.Find("div", "class", "updates").FullText()
			text = pattern.ReplaceAllString(text, "\n")

			small := doc.Find("small")
			date := time.Now().String()
			if small.Error == nil {
				date = small.Text()
			}
			events = append(
				events,
				crawler.Event{
					Crawler:  c.GetCrawlerName(),
					Job:      c.GetJobName(),
					UserName: c.channel,
					UID:      date,
					Name:     date,
					Message:  fmt.Sprintf("%s%s<%s|Confluent Cloud Status>", a.FullText(), text, url),
				},
			)
		}
	}

	return events, nil
}

func (c Crawler) containsKeyword(s string) bool {
	for _, keyword := range c.keywords {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}

func (c Crawler) CrawlRelease() ([]crawler.Event, error) {
	var jsonBody string
	var dtos []DTO

	url := "https://docs.confluent.io/cloud/current/release-notes/index.html"

	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*2),
		chromedp.EvaluateAsDevTools(
			`
			function map_object(div) {
				return {
					"date": div.querySelector('h2').innerText,
					"content": div.innerText,
				};
			}

			function crawl() {
				var records = [];
				var divs = document.querySelectorAll('div#ccloud-release-notes > div.section')
				for (var i = 0; i < Math.min(divs.length, 5); i++) { // 최근 5개만 확인
					records.push(map_object(divs[i]));
				}
				
				return JSON.stringify(records.reverse()); // 과거 릴리즈가 인덱스 앞쪽으로 하기위해 reverse
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

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel, url)
}

func NewCrawler(logger *zap.Logger, chromectx context.Context, channel string, jobName string, keywords []string) *Crawler {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,

		channel:  channel,
		jobName:  jobName,
		keywords: keywords,
	}
}
