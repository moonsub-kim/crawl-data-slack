package slackclient

import (
	"fmt"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/logger"
	"github.com/slack-go/slack"
)

const CHANNEL_ID = "G015JFASK7Z"

var messageBuilder = map[string]map[string](func(crawler.Notification) string){
	"groupware": map[string]func(crawler.Notification) string{
		"declined_payments": func(n crawler.Notification) string {
			return fmt.Sprintf("<@%s> 문서가 반려되었습니다. 그룹웨어에서 확인해주세요.", n.User.ID)
		},
	},
}

type Client struct {
	logger logger.Logger
	api    *slack.Client
	mapper mapper
}

func (c Client) Notify(n crawler.Notification) error {
	message := messageBuilder[n.Event.Crawler][n.Event.Job](n)
	_, _, _, err := c.api.SendMessage(CHANNEL_ID, slack.MsgOptionText(message, false))
	return err
}

func (c Client) GetUsers() ([]crawler.User, error) {
	users, err := c.api.GetUsers()
	if err != nil {
		return nil, err
	}

	return c.mapper.mapSlackUsersToUsers(users), nil
}

func NewClient(logger logger.Logger, client *slack.Client) *Client {
	return &Client{
		logger: logger,
		api:    client,
	}
}
