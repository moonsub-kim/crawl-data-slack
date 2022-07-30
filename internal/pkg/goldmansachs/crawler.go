package goldmansachs

import (
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel    string
	recentDays int
}

const URL string = "https://developer.gs.com/blog/posts"
const BASE_URL string = "https://developer.gs.com"
const DATE_PARSE_FORMAT string = "January 2 2006"

func (c Crawler) GetCrawlerName() string { return "goldmansachs" }
func (c Crawler) GetJobName() string     { return "post" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}

	t := time.Now().Add(time.Duration(-c.recentDays) * time.Hour * 24)

	doc := soup.HTMLParse(res)
	as := doc.Find("div", "class", "gs-uitk-c-1c4ow0d").FindAll("a")
	var dtos []DTO
	for _, a := range as {
		spans := a.FindAll("span")
		dateStr := strings.ReplaceAll(spans[0].Text(), ",", "")
		name := spans[1].Text()
		url := BASE_URL + a.Attrs()["href"]

		date, err := time.Parse(DATE_PARSE_FORMAT, dateStr)
		if err != nil {
			return nil, err
		}

		if date.Before(t) {
			continue
		}

		dtos = append(dtos, DTO{
			ID:   strings.TrimSpace(name),
			Date: strings.TrimSpace(dateStr),
			Name: strings.TrimSpace(name),
			URL:  strings.TrimSpace(url),
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string, recentDays int) *Crawler {
	return &Crawler{
		logger:     logger,
		channel:    channel,
		recentDays: recentDays,
	}
}
