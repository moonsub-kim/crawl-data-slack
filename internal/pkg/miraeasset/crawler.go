package miraeasset

import (
	"fmt"
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

const CONTENT_URL_TEMPLATE string = "https://securities.miraeasset.com/bbs/board/message/view.do?messageId=%s&messageNumber=%s&categoryId=1521"
const PDF_URL_TEMPLATE string = "https://securities.miraeasset.com/bbs/download/%s.pdf?attachmentId=%s"

func (c Crawler) GetCrawlerName() string { return "ipo" }
func (c Crawler) GetJobName() string     { return "ipo" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var dtos []DTO

	res, err := soup.Get("https://securities.miraeasset.com/bbs/board/message/list.do?categoryId=1521")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(res)

	trs := doc.Find("table", "class", "bbs_linetype2").Find("tbody").FindAll("tr")
	for _, tr := range trs {
		tds := tr.FindAll("td")
		date := strings.TrimSpace(tds[0].Text())
		a := tds[1].Find("div", "class", "subject").Find("a")

		title := strings.ReplaceAll(strings.TrimSpace(a.FullText()), "\n", " ")
		parsed := strings.Split(a.Attrs()["href"], "'")
		id := parsed[1]
		number := parsed[3]
		contentURL := fmt.Sprintf(CONTENT_URL_TEMPLATE, id, number)
		pdfURL := ""
		if tds[2].Find("p").Error == nil {
			pdfURL = fmt.Sprintf(PDF_URL_TEMPLATE, id, id)
		} // pdfURL이 있는경우만

		content, err := c.parseContent(contentURL)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, DTO{
			ID:      id,
			Date:    date,
			Title:   title,
			pdfURL:  pdfURL,
			URL:     contentURL,
			Content: content,
		})
	}

	return c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel), nil
}

func (c Crawler) parseContent(contentURL string) (string, error) {
	res, err := soup.Get(contentURL)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(res)
	div := doc.Find("div", "id", "messageContentsDiv")

	table := div.Find("table")
	var tabletxt string
	if table.Error == nil {
		// div.Pointer.RemoveChild(table.Pointer)
		tabletxt = c.buildTable(table)
	}

	m := regexp.MustCompile(`\n+`)
	text := m.ReplaceAllString("> "+div.FullText(), "\n> ")
	if tabletxt != "" {
		text += "```\n" + tabletxt + "```"
	}

	return text, nil
}

func (c Crawler) buildTable(table soup.Root) string {
	grid := [][]string{}
	fmt.Println(table.HTML())

	ths := table.FindAll("th")
	rows := table.FindAll("tr")

	grid = append(grid, []string{})
	for _, th := range ths {
		grid[len(grid)-1] = append(grid[len(grid)-1], th.Text())
	}

	for _, row := range rows {
		tds := row.FindAll("td")
		if len(tds) == len(ths) {
			grid = append(grid, []string{})
		} else {
			continue
		}
		for _, td := range tds {
			grid[len(grid)-1] = append(grid[len(grid)-1], td.Text())
			fmt.Printf("append %s", td.Text())
		}
	}

	fmt.Println(grid)

	for j := 0; j < len(ths); j++ { // col
		// get max
		maxLen := 0
		for i := 0; i < len(grid); i++ { // row
			if maxLen < len(grid[i][j]) {
				maxLen = len(grid[i][j])
			}
		}
		maxLen += 1

		// padding
		for i := 0; i < len(grid); i++ {
			grid[i][j] = fmt.Sprintf(fmt.Sprintf("%%-%ds", maxLen-len(grid[i][j])), grid[i][j])
		}
	}

	txt := ""
	for i := 0; i < len(grid); i++ {
		s := strings.Join(grid[i], " | ")
		txt += s + "\n"
	}

	return txt
}

func NewCrawler(logger *zap.Logger, channel string) *Crawler {
	return &Crawler{
		logger:  logger,
		channel: channel,
	}
}
