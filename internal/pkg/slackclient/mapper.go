package slackclient

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/slack-go/slack"
)

type mapper struct {
}

func (mapper) mapSlackUserToUser(user slack.User) crawler.User {
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
