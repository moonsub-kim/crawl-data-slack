package slackclient

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/slack-go/slack"
)

type mapper struct {
}

func (m mapper) mapSlackUserToUser(user slack.User) crawler.User {
	return crawler.User{
		ID:   user.ID,
		Name: user.Name,
	}
}

func (m mapper) mapSlackUsersToUsers(slackUsers []slack.User) []crawler.User {
	var users []crawler.User
	for _, u := range slackUsers {
		users = append(users, m.mapSlackUserToUser(u))
	}
	return users
}

func (m mapper) mapSlackChannelsToUsers(slackChannels []slack.Channel) []crawler.User {
	var users []crawler.User
	for _, c := range slackChannels {
		users = append(users, m.mapSlackChannelToUser(c))
	}
	return users
}

func (m mapper) mapSlackChannelToUser(c slack.Channel) crawler.User {
	return crawler.User{
		ID:   c.ID,
		Name: c.Name,
	}
}
