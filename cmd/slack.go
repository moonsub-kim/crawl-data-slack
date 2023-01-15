package main

import (
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	slackListConversationsArgChannelID = "channel-id"
	slackListConversationsArgFromDate  = "from-date"
	slackListConversationsArgToDate    = "to-date"

	commandArchivePosts *cli.Command = &cli.Command{
		Name: "archive-posts",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: slackListConversationsArgChannelID},
			&cli.TimestampFlag{Name: slackListConversationsArgFromDate, Layout: dateLayout},
			&cli.TimestampFlag{Name: slackListConversationsArgToDate, Layout: dateLayout},
		},
		Action: RunSlack(
			func(ctx *cli.Context, logger *zap.Logger, client *slackclient.Client) error {
				channel := ctx.String(slackListConversationsArgChannelID)
				from := ctx.Timestamp(slackListConversationsArgFromDate).Add(-time.Hour * 9) // UTC->KST
				to := ctx.Timestamp(slackListConversationsArgToDate).Add(time.Hour * 24).Add(-time.Hour * 9)

				_, err := client.ArchivePosts(crawler.Channel{ID: channel}, from, to)
				return err
			},
		),
	}
)
