package kcif

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	channel      string
	eventBuilder eventBuilder
	query        string
}

const URL string = "https://www.kcif.or.kr/front/board/boardList.do?intSection1=2"
const REFERRER string = "https://www.kcif.or.kr/front/board/boardList.do?intSection1=2"
const PDF_URL_BASE string = "https://www.kcif.or.kr/front/board/fileDownLoad.do?board_id=%s&fileGb=%s"
const VIEW_URL string = "https://www.kcif.or.kr/front/board/boardView.do"

func (c Crawler) GetCrawlerName() string { return "kcif" }
func (c Crawler) GetJobName() string     { return "report" }

func (c Crawler) request(method string, url string, referrer string, reqBody io.Reader, header map[string]string) (string, error) {
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Referer", referrer)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7,ru;q=0.6")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", `Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "macOS")

	for k, v := range header {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO

	res, err := c.request(http.MethodGet, URL, REFERRER, nil, nil)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)
	trs := doc.Find("tbody").FindAll("tr")

	postHeader := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Origin":       "https://www.kcif.or.kr",
	}
	for _, tr := range trs {
		tds := tr.FindAll("td")
		// 공개된 문서가 아니면 제외
		if tds[0].Find("img", "class", "open").Error != nil {
			continue
		}

		date := tds[2].Text()

		pdfFuncStr := tds[4].Find("a").Attrs()["onclick"] // `download('50261', '1')`
		pattern := regexp.MustCompile(`'(?P<id>\d+)', ?'(?P<index>\d+)'`)
		matches := pattern.FindStringSubmatch(pdfFuncStr)

		id := matches[1]
		fileGb := matches[2]
		pdfURL := fmt.Sprintf(PDF_URL_BASE, id, fileGb)

		querystr := fmt.Sprintf("intReportID=%s&currentPage=1&rowsPerPage=15&intSection1=2&intSection2=5&intBoardID=5&regular=&AnalysisBrief=&intperiod1=&intperiod2=&orderValue=&s_title=true&s_word=", id)
		res, err = c.request(http.MethodPost, VIEW_URL, VIEW_URL, strings.NewReader(querystr), postHeader)
		if err != nil {
			c.logger.Error(
				"Failed to request",
				zap.Error(err),
			)
			continue
		}

		doc = soup.HTMLParse(res)
		title := doc.Find("td", "id", "title").Text()
		content := doc.Find("td", "id", "contents").Find("p").Find("span").FullText()

		dtos = append(dtos, DTO{
			ID:      strings.TrimSpace(id),
			Date:    strings.TrimSpace(date),
			Title:   strings.TrimSpace(title),
			pdfURL:  strings.TrimSpace(pdfURL),
			Content: strings.TrimSpace(content),
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
