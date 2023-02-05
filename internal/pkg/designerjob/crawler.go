package designerjob

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger   *zap.Logger
	channel  string
	query    string
	excludes []string
}

func (c Crawler) GetCrawlerName() string { return "desginer-job" }
func (c Crawler) GetJobName() string     { return "open-position" }

func (c Crawler) isExcludes(title string) bool {
	for _, s := range c.excludes {
		if strings.Contains(title, s) {
			return true
		}
	}
	return false
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	url := fmt.Sprintf(
		"https://www.designerjob.co.kr/search/total-search?%s=%s&%s=%s",
		url.QueryEscape("param[search]"),
		url.QueryEscape(c.query),
		url.QueryEscape("param[cmd]"),
		"recurit",
	)
	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	var events []crawler.Event
	p := regexp.MustCompile("[|\n\t]")
	doc := soup.HTMLParse(res)
	lis := doc.Find("div", "class", "con_list").Find("ul").FindAll("li")
	for _, li := range lis {
		text := p.ReplaceAllString(li.FullText(), "")
		if c.isExcludes(strings.ToLower(text)) {
			continue
		}

		url := li.Find("a", "class", "subject").Attrs()["href"]

		events = append(
			events,
			crawler.Event{
				Crawler:   c.GetCrawlerName(),
				Job:       c.GetJobName(),
				UserName:  c.channel,
				UID:       url,
				Name:      text,
				EventTime: time.Now(), // TODO exact event time
				Message:   fmt.Sprintf("<%s|%s>", url, text),
			},
		)
	}

	return events, nil
}

func NewCrawler(logger *zap.Logger, channel string, query string, excludes []string) *Crawler {
	return &Crawler{
		logger:   logger,
		channel:  channel,
		query:    url.QueryEscape(query),
		excludes: excludes,
	}
}
