package interpark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
	date    string
}

const URL string = "https://api-ticketfront.interpark.com/v1/goods/22015433/playSeq/PlayDate/%s/REMAINSEAT"

func (c Crawler) GetCrawlerName() string { return "interpark" }
func (c Crawler) GetJobName() string     { return "ticket" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	url := fmt.Sprintf(URL, c.date)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
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

	return c.eventBuilder.buildEvents(resDTO, c.date, c.GetCrawlerName(), c.GetJobName(), c.channel)
}

func NewCrawler(logger *zap.Logger, channel string, date string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
		date:    date,
	}
}
