package financialreport

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder

	channel string
}

const URL string = "https://www.kcif.or.kr/front/board/boardList.do?intSection1=2"
const CONTENTS int = 15
const SCAN_PAGES int = 2

func (c Crawler) GetCrawlerName() string { return "kcif-report" }
func (c Crawler) GetJobName() string     { return "kcif" }

// just crawl first page
func (c Crawler) Crawl() ([]crawler.Event, error) {
	dtos := make([]DTO, CONTENTS*SCAN_PAGES)
	var actions []chromedp.Action
	for i := 0; i < SCAN_PAGES; i++ {
		for j := 0; j < CONTENTS; j++ {
			actions = append(actions, c.crawlActions(URL, j, &dtos[i*CONTENTS+j])...)
		}
		nextPage := i + 1
		actions = append(actions, c.nextPageActions(nextPage)...)
	}

	err := chromedp.Run(
		c.ctx,
		actions...,
	)
	if err != nil {
		c.logger.Error("run error", zap.Error(err))
		return nil, err
	}

	targets, err := chromedp.Targets(c.ctx)
	if err != nil {
		c.logger.Error("retrieving targets error", zap.Error(err))
		return nil, err
	}
	c.logger.Info("targets", zap.Any("targets", targets))

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.channel)
	if err != nil {
		return nil, err
	}

	c.logger.Info(
		"events",
		zap.Any("events", events),
	)

	return events, nil
}


// 실제 크롤링 수행
func (c Crawler) crawlActions(link string, j int, dto *DTO) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(link),
		chromedp.Sleep(time.Second * 2),
		chromedp.Evaluate(
			fmt.Sprintf(
				`
				function pdfTab(tdPDF) {
					tdPDF.querySelector('a').click();
				}
				
				function getText(tdTitle) {
					tdTitle.querySelector('a').click(); // 창 열기
					text = document.querySelector('div.ng-binding').innerText;
					document.querySelectorAll('div.right > a.cursor')[1].click(); // 창 닫기
					return text;
				}

				function mapObject(tr) {
					tds = tr.querySelectorAll('td');
					pdfTab(tds[2]);

					return {
						'date': tds[0].innerText, // YYYY/MM/DD format
						'title': tds[1].innerText,
						'text': getText(tds[1]),
						'company': tds[3].innerText,
					}
				}
				
				function main() {
					var trs = document.querySelectorAll('table.report-table > tbody > tr');
					return mapObject(trs[%d]);
				}

				main()
				`,
				j,
			),
			dto,
		),
		chromedp.ActionFunc(func(ctx context.Context) error {
			targets, err := chromedp.Targets(ctx)
			if err != nil {
				c.logger.Info("Targets", zap.Error((err)))
				return err
			}

			for _, t := range targets {
				if t.URL == "about:blank" || t.URL == URL {
					continue
				}
				if !t.Attached {
					dto.PDFURL = t.URL
					newCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(t.TargetID))
					if err := chromedp.Run(newCtx, chromedp.Sleep(time.Millisecond*100)); err != nil {
						c.logger.Error("run err", zap.Error(err))
					}
				}
			}
			return nil
		}),
	}
}

func NewCrawler(logger *zap.Logger, chromectx context.Context, channel string) (*Crawler, error) {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,

		channel: channel,
	}, nil
}
