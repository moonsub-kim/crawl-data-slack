package main

import (
	"fmt"
	"os"
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
	slackListConversationsArgFilter    = "filter"

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
			negatives := map[string]slackclient.ArchiveFilter{"no-link": slackclient.NoLinkFilter{}}
			for _, f := range ctx.StringSlice(slackListConversationsArgFilter) {
				filter, ok := negatives[f]
				if !ok {
					return fmt.Errorf("no adequate filter for %s", f)
				}
				negativeFilters = append(negativeFilters, filter)
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
