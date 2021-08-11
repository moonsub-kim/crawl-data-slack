package slackclient

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	api    *slack.Client
	mapper mapper
}

func (c Client) Notify(n crawler.Notification) error {
	_, _, err := c.api.PostMessage(n.User.ID, slack.MsgOptionText(n.Event.Message, false))
	c.logger.Info(
		"notify",
		zap.Any("notification", n),
		zap.Error(err),
	)
	return err
}

func (c Client) GetUsers() ([]crawler.User, error) {
	users, err := c.api.GetUsers()
	if err != nil {
		c.logger.Error("getUsers", zap.Error(err))
		return nil, err
	}

	var activeUsers []slack.User
	for _, u := range users {
		if u.Deleted || u.IsBot || u.IsRestricted {
			continue
		}

		activeUsers = append(activeUsers, u)
	}

	return c.mapper.mapSlackUsersToUsers(activeUsers), nil
}

func NewClient(logger *zap.Logger, client *slack.Client) *Client {
	return &Client{
		logger: logger,
		api:    client,
	}
}
