package groupware

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger *zap.Logger
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	// goto:: https://gr.buzzvil.com/gw/uat/uia/egovLoginUsr.do
	// #userId -> id넣기
	// #userPw -> pw넣기
	// .log_btn >> input -> 클릭
	// ((TEST)) #2002010000_anchor -> 클릭,, ((PROD)) #2002070000_anchor -> 클릭
	// ((TEST)) #from_date -> 2달전으로, ((PROD)) #from_date -> value 넣기 전날로 .
	// #searchBtn -> 클릭
	// .grid-content >> table >> tbody >> tr 한개씩읽기
	return nil, nil
}

func NewCrawler(logger *zap.Logger) *Crawler {
	return &Crawler{
		logger: logger,
	}
}
