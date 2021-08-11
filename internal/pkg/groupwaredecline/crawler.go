package groupwaredecline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder
	id           string
	pw           string
}

func (c Crawler) GetCrawlerName() string { return "groupware" }
func (c Crawler) GetJobName() string     { return "declined_payments" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var jsonBody string
	var dtos []DTO

	yesterday := time.Now().Add(time.Hour*9).AddDate(0, -2, -1).Format("2006-01-02") // UTCNOW -> KST -> yesterday -> formating
	c.logger.Info("yesterday", zap.String("yesterday", yesterday))
	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate("https://gr.buzzvil.com/gw/uat/uia/egovLoginUsr.do"),

		// 로그인페이지: 로그인
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`document.getElementById('userId').value = '%s'; document.getElementById('userPw').value = '%s'; actionLogin();`,
				c.id,
				c.pw,
			),
			nil,
		),
		chromedp.Sleep(time.Second*2),

		// 반려함, (test용 결재요청함 2002010000)
		chromedp.Navigate(`https://gr.buzzvil.com/eap/ea/eadoc/EaDocList.do?menu_no=2002070000`),

		// 전자 결재 - 반려함: 일자 조정 버튼
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`document.getElementById('from_date').value = '%s'; document.getElementById('searchBtn').click();`,
				yesterday,
			),
			nil,
		),
		chromedp.Sleep(time.Second*2),

		// 문서 파싱
		chromedp.Evaluate(
			`
			function map_object(arr) {
				const indexMap = {2: "id", 3: "doc_name", 4: "request_date", 9: "drafter", 10: "status"};
				var keys = Object.keys(indexMap);
				var obj = {};

				for (var i = 0; i < keys.length; i++) {
					k = indexMap[keys[i]];
					obj[k] = arr[keys[i]];
				}

				return obj;
			}

			function crawl() {
				var trs = document.body.querySelectorAll('div.grid-content > table > tbody > tr');
				var records = [];
				for (var i = 0; i < trs.length; i++) {
					var arr = [];
					var tds = trs[i].getElementsByTagName('td');
					for (var j = 0; j < tds.length; j++) {
						arr.push(tds[j].innerText);
					}
					console.log(arr)
					records.push(map_object(arr));
				}

				return JSON.stringify(records);
			}
			crawl();
			`,
			&jsonBody,
		),
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonBody), &dtos)
	if err != nil {
		return nil, err
	}

	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName())
	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewCrawler(logger *zap.Logger, chromectx context.Context, id string, pw string) *Crawler {
	return &Crawler{
		logger: logger,
		ctx:    chromectx,
		id:     id,
		pw:     pw,
	}
}
