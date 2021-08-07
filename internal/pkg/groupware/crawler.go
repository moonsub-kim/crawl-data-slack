package groupware

import "github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"

type Crawler struct {
}

func (c Crawler) Crawl() ([]crawler.Event, error) {
	return nil, nil
}

func NewCrawler() *Crawler {
	return &Crawler{}
}
