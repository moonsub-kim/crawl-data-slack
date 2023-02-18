package naverd2

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://d2.naver.com/api/v1/contents?categoryId=2&page=0&size=2"
const BASE_URL string = "https://d2.naver.com"

func (c Crawler) GetCrawlerName() string { return "naver" }
func (c Crawler) GetJobName() string     { return "d2" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resDTO Response
	err = json.Unmarshal(body, &resDTO)
	if err != nil {
		return nil, err
	}

	return c.eventBuilder.buildEvents(resDTO.Contents, c.GetCrawlerName(), c.GetJobName(), c.channel, BASE_URL), nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
