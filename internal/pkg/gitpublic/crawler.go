package gitpublic

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel      string
	organization string
}

func (c Crawler) GetCrawlerName() string { return "git" }
func (c Crawler) GetJobName() string     { return "public_repo" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	url := fmt.Sprintf(`https://github.com/orgs/%s/repositories`, c.organization)

	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	lis := doc.Find("div", "class", "org-repos").Find("div", "class", "Box").Find("ul").FindAll("li")
	var dtos []DTO
	for _, li := range lis {
		a := li.Find("a", "class", "d-inline-block")
		name := a.Text()
		url := "https://github.com" + a.Attrs()["href"]

		dtos = append(dtos, DTO{
			ID:   strings.TrimSpace(name),
			Name: strings.TrimSpace(name),
			URL:  strings.TrimSpace(url),
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string, organization string) *Crawler {
	return &Crawler{
		logger:       logger,
		channel:      channel,
		organization: organization,
	}
}
