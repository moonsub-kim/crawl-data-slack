package main

import (
	"context"
	"os"
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
	slackArchiveArgChannel = "channel"
	// slackArchiveArgFromDate = "from-date"
	// slackArchiveArgToDate   = "to-date"

	commandArchive *cli.Command = &cli.Command{
		Name: "archive",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: slackArchiveArgChannel, Required: true},
			// &cli.TimestampFlag{Name: slackArchiveArgFromDate, Layout: dateLayout},
			// &cli.TimestampFlag{Name: slackArchiveArgToDate, Layout: dateLayout},
		},
		Action: func(ctx *cli.Context) error {
			slackBotToken := os.Getenv(envSlackBotToken)
			githubToken := os.Getenv(envGithubToken)

			logger := zapLogger()

			kv := map[string]string{}
			for _, k := range ctx.FlagNames() {
				kv[k] = ctx.String(k)
			}
			logger.Info(
				"flags",
				zap.Any("flags", kv),
			)

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
			weekday := time.Duration(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			year, month, day := now.Date()
			currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			dateFrom := currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
			dateFrom = dateFrom.Add(-time.Hour * 24 * 7)
			logger.Info("date", zap.Any("date", dateFrom))

			err = u.Archive(channel, dateFrom, now)
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
