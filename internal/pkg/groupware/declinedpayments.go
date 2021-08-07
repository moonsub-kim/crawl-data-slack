package groupware

import "github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"

type DeclinedPaymentsCrawler struct {
	groupwareCrawler
}

func (c DeclinedPaymentsCrawler) Crawl() ([]crawler.Event, error) {
	return nil, nil
}

func NewDeclinedPaymentsCrawler() *DeclinedPaymentsCrawler {
	return &DeclinedPaymentsCrawler{
		groupwareCrawler: groupwareCrawler{},
	}
}
