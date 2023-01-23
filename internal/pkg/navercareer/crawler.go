package navercareer

import (
	"io"
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/google/go-querystring/query"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder
	channel      string
	query        string
	excepts      []string
}

const URL string = "https://s.search.naver.com/p/career/search.naver"

func (c Crawler) GetCrawlerName() string { return "naver" }
func (c Crawler) GetJobName() string     { return "career" }

func (c Crawler) includeExcepts(title string) bool {
	for _, s := range c.excepts {
		if strings.Contains(title, s) {
			return true
		}
	}
	return false
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	reqBody, err := newPageRequest(c.query)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, URL+"?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	s := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				string(b), `\"`, `"`,
			), `({ results:[ "`, "",
		), `" ] });`, "",
	)
	splitted := strings.Split(s, `", "`)

	c.logger.Info("splitted", zap.Any("splitted", splitted))
	var dtos []DTO
	for _, s := range splitted {
		doc := soup.HTMLParse(s)

		if doc.Find("div").Error != nil {
			continue
		}

		div_title_area := doc.Find("div", "class", "title_area")
		div_info_area := doc.Find("div", "class", "info_area")

		title := strings.TrimSpace(div_title_area.FullText())
		if c.includeExcepts(title) {
			continue
		}

		dtos = append(dtos, DTO{
			Title: title,
			Info:  strings.TrimSpace(div_info_area.FullText()),
			URL:   div_title_area.Find("a", "class", "title").Attrs()["href"],
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func NewCrawler(logger *zap.Logger, channel string, query string, excepts []string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
		query:   query,
		excepts: excepts,
	}
}
