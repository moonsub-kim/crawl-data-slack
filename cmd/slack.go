package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
)

var (
	slackListConversationsArgChannelID = "channel-id"
	slackListConversationsArgFromDate  = "from-date"
	slackListConversationsArgToDate    = "to-date"
	slackListConversationsArgFilter    = "filter" // --filter no-link --filter exclude-emoji:emoji_name

	commandListConversations *cli.Command = &cli.Command{
		Name: "list-conversations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: slackListConversationsArgChannelID},
			&cli.TimestampFlag{Name: slackListConversationsArgFromDate, Layout: dateLayout},
			&cli.TimestampFlag{Name: slackListConversationsArgToDate, Layout: dateLayout},
			&cli.StringSliceFlag{Name: slackListConversationsArgFilter},
		},
		Action: func(ctx *cli.Context) error {
			slackBotToken := os.Getenv(envSlackBotToken)

			logger := zapLogger(ctx)

			var negativeFilters []slackclient.ArchiveFilter
			for _, f := range ctx.StringSlice(slackListConversationsArgFilter) {
				splitted := strings.Split(f, ":")

				if splitted[0] == "no-link" {
					negativeFilters = append(negativeFilters, slackclient.NoLinkFilter{})
				} else if splitted[0] == "exclude-emoji" {
					if len(splitted) != 2 {
						return fmt.Errorf("exclude-emoji filter requires only one arguement")
					}
					negativeFilters = append(negativeFilters, slackclient.NewExcludeEmojiFilter(splitted[1]))
				} else {
					return fmt.Errorf("no adequate filter for %s", splitted[0])
				}
			}

			client := slackclient.NewClient(
				logger,
				slack.New(slackBotToken),
				slackBotToken,
				negativeFilters,
			)

			channel := ctx.String(slackListConversationsArgChannelID)
			from := ctx.Timestamp(slackListConversationsArgFromDate).Add(-time.Hour * 9) // UTC->KST
			to := ctx.Timestamp(slackListConversationsArgToDate).Add(-time.Hour * 9)

			_, err := client.ArchivePosts(crawler.Channel{ID: channel}, from, to)
			return err
		},
	}
)
