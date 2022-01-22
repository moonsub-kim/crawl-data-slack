package slackclient

import (
	"time"

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
	time.Sleep(time.Second)
	_, _, err := c.api.PostMessage(n.User.ID, slack.MsgOptionText(n.Event.Message, false))
	c.logger.Info(
		"notify",
		zap.Any("notification", n),
		zap.Any("err", err),
	)
	return err
}

func (c Client) GetChannels() ([]crawler.Channel, error) {
	slackUsers, err := c.api.GetUsers()
	if err != nil {
		c.logger.Error("getUsers", zap.Error(err))
		return nil, err
	}

	var activeUsers []slack.User
	for _, u := range slackUsers {
		if u.Deleted || u.IsBot || u.IsRestricted {
			continue
		}

		activeUsers = append(activeUsers, u)
	}

	users := c.mapper.mapSlackUsersToUsers(activeUsers)

	nextCursor := ""
	for {
		var slackChannels []slack.Channel
		param := slack.GetConversationsParameters{Cursor: nextCursor, ExcludeArchived: true}
		slackChannels, nextCursor, err = c.api.GetConversations(&param)
		if err != nil {
			c.logger.Error("getConversations", zap.Error(err))
			return nil, err
		} else if nextCursor == "" {
			break
		}
		channels := c.mapper.mapSlackChannelsToUsers(slackChannels)
		users = append(users, channels...)
		c.logger.Info(
			"GetConversations",
			zap.Any("channels", channels),
			zap.Any("nextCursor", nextCursor),
		)
		time.Sleep(time.Second * 3)
	}

	return users, nil
}

func NewClient(logger *zap.Logger, client *slack.Client) *Client {
	return &Client{
		logger: logger,
		api:    client,
	}
}
