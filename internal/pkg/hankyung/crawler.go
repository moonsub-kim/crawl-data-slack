package hankyung

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger  *zap.Logger
	channel string
}

const URL string = "https://www.hankyung.com/globalmarket/news/wallstreet-now"

func (c Crawler) GetCrawlerName() string { return "hankyung" }
func (c Crawler) GetJobName() string     { return "wallstreetnow" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := soup.Get(URL)
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	news := doc.Find("div", "class", "list_thumb_rowtype")
	a := news.Find("h3").Find("a")
	// datetime := news.Find("span", "class", "time")

	title := strings.ReplaceAll(a.FullText(), "회원전용 ", "")
	url := a.Attrs()["href"]
	// createdAt := datetime.FullText()

	res, err = soup.Get(url)
	if err != nil {
		return nil, err
	}
	doc = soup.HTMLParse(res)
	figureReplacer := regexp.MustCompile(`<img src="(.+?)".+>`)
	iframeRemover := regexp.MustCompile(`<iframe.+?</iframe>`)
	tagRemover := regexp.MustCompile(`<.+?>`)
	consecutiveNewlineRemover := regexp.MustCompile(`(\n\s*)+`)
	body := doc.Find("div", "class", "article-body").HTML()
	body = figureReplacer.ReplaceAllString(body, "$1")
	body = iframeRemover.ReplaceAllString(body, "")
	body = tagRemover.ReplaceAllString(body, "")
	body = consecutiveNewlineRemover.ReplaceAllString(body, "\n")
	body = html.UnescapeString(body)

	return []crawler.Event{
		{
			Crawler:  c.GetCrawlerName(),
			Job:      c.GetJobName(),
			UserName: c.channel,
			UID:      title,
			Name:     title,
			Message:  fmt.Sprintf("%s\n<%s|URL>", body, url),
		},
	}, nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
