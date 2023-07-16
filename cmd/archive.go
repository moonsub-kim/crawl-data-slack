package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/githubclient"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	commandSyncLabel *cli.Command = &cli.Command{
		Name: "sync-label",
		Action: func(ctx *cli.Context) error {
			slackBotToken := os.Getenv(envSlackBotToken)
			githubToken := os.Getenv(envGithubToken)

			logger := zapLogger(ctx)

			db, err := openDB(logger)
			if err != nil {
				logger.Info("openDB", zap.Error(err))
				return err
			}

			u := crawler.NewUseCase(
				logger,
				repository.NewRepository(logger, db),
				nil,
				slackclient.NewClient(
					logger,
					slack.New(slackBotToken),
					slackBotToken,
					nil,
				),
				githubclient.NewClient(
					logger,
					github.NewClient(
						oauth2.NewClient(
							context.Background(),
							oauth2.StaticTokenSource(
								&oauth2.Token{AccessToken: githubToken},
							),
						),
					),
					ctx.String("owner"),
					ctx.String("repo"),
				),
			)

			return u.SyncLabel()
		},
	}
)

var (
	slackArchiveArgChannel  = "channel"
	slackArchiveArgFromDate = "from-date"
	slackArchiveArgToDate   = "to-date"
	slackArchiveArgFilter   = "filter" // --filter no-link --filter exclude-emoji:emoji_name

	commandArchive *cli.Command = &cli.Command{
		Name: "archive",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: slackArchiveArgChannel, Required: true},
			&cli.TimestampFlag{Name: slackArchiveArgFromDate, Layout: dateLayout},
			&cli.TimestampFlag{Name: slackArchiveArgToDate, Layout: dateLayout},
			&cli.StringSliceFlag{Name: slackArchiveArgFilter},
		},
		Action: func(ctx *cli.Context) error {
			slackBotToken := os.Getenv(envSlackBotToken)
			githubToken := os.Getenv(envGithubToken)

			logger := zapLogger(ctx)

			db, err := openDB(logger)
			if err != nil {
				logger.Info("openDB", zap.Error(err))
				return err
			}

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

			u := crawler.NewUseCase(
				logger,
				repository.NewRepository(logger, db),
				nil,
				slackclient.NewClient(
					logger,
					slack.New(slackBotToken),
					slackBotToken,
					negativeFilters,
				),
				githubclient.NewClient(
					logger,
					github.NewClient(
						oauth2.NewClient(
							context.Background(),
							oauth2.StaticTokenSource(
								&oauth2.Token{AccessToken: githubToken},
							),
						),
					),
					ctx.String("owner"),
					ctx.String("repo"),
				),
			)

			channel := ctx.String(slackArchiveArgChannel)
			now := time.Now()
			toDate := now
			weekday := time.Duration(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			year, month, day := now.Date()
			currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			fromDate := currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
			fromDate = fromDate.Add(-time.Hour * 24 * 7)

			if ctx.Timestamp(slackArchiveArgFromDate) != nil {
				fromDate = ctx.Timestamp(slackArchiveArgFromDate).Add(-time.Hour * 9)
			}
			if ctx.Timestamp(slackArchiveArgToDate) != nil {
				toDate = ctx.Timestamp(slackArchiveArgToDate).Add(-time.Hour * 9)
			}

			logger.Info("date", zap.Time("from_date", fromDate), zap.Time("to_date", toDate))
			err = u.Archive(channel, fromDate, toDate)
			if err != nil {
				logger.Error(
					"Archive Error",
					zap.Error(err),
				)
				return err
			}

			logger.Info("Succeeded")
			return nil
		},
	}
)
