package lottecinema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	eventBuilder eventBuilder

	channel string
	date    string
}

const URL string = "https://www.lottecinema.co.kr/LCWS/Ticketing/TicketingData.aspx"
const DATA_FORMAT string = `------WebKitFormBoundaryziISgzfxg73lJkgP
Content-Disposition: form-data; name="paramList"

{"MethodName":"GetPlaySequence","channelType":"HO","osType":"W","osVersion":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36","playDate":"%s","cinemaID":"1|0001|1016","representationMovieCode":"18755"}
------WebKitFormBoundaryziISgzfxg73lJkgP--
`

func (c Crawler) GetCrawlerName() string { return "lottecinema" }
func (c Crawler) GetJobName() string     { return "movie" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	reqBody := fmt.Sprintf(DATA_FORMAT, c.date)
	req, err := http.NewRequest(http.MethodGet, URL, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryziISgzfxg73lJkgP")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resDTO Response
	err = json.Unmarshal(body, &resDTO)
	if err != nil {
		return nil, err
	}

	return c.eventBuilder.buildEvents(resDTO.PlaySeqs.Items, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func NewCrawler(logger *zap.Logger, channel string, date string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
		date:    date,
	}
}
