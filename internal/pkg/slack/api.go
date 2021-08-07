package slack

import "github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"

type Service struct {
}

func (s Service) Notify(n crawler.Notification) error {
	return nil
}

func (s Service) GetUsers() ([]crawler.User, error) {
	return nil, nil
}

func NewService() *Service {
	return &Service{}
}
