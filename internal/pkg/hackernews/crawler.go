package hackernews

import (
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
	filters      []Filter
}

func (c Crawler) GetCrawlerName() string { return "hacker-news" }
func (c Crawler) GetJobName() string     { return "news" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO

	res, err := soup.Get("https://news.ycombinator.com/news")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	contents := doc.Find("table", "class", "itemlist").Find("tbody").FindAll("tr")
	for i := 0; i < len(contents); i += 3 {
		athing := contents[i]
		metdata := contents[i+1]
		if athing.Attrs()["class"] == "morespace" {
			break
		}

		id := strings.TrimSpace(athing.Attrs()["id"])
		a := athing.Find("a", "class", "storylink")
		href := strings.TrimSpace(a.Attrs()["href"])
		if strings.HasPrefix(href, "item?id=") {
			href = "https://news.ycombinator.com/" + href
		}
		subtext := metdata.Find("td", "class", "subtext")

		dtos = append(dtos, DTO{
			ID:         id,
			URL:        href,
			CommentURL: "https://news.ycombinator.com/item?id=" + id,
			Title:      strings.TrimSpace(a.Text()),
			SubText:    strings.TrimSpace(subtext.FullText()),
		})
	}

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel, c.filters)
	if err != nil {
		return nil, err
	}

	c.logger.Info(
		"crawler",
		zap.Any("dto", dtos),
		zap.Any("events", events),
	)
	return events, nil
}

func NewCrawler(logger *zap.Logger, channel string, pointThreshold int) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
		filters: []Filter{
			&adFilter{},
			&ageFilter{},
			&pointFilter{threshold: pointThreshold},
			&commentFilter{},
		},
	}
}
