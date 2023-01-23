package main

import (
	"context"
	"fmt"
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
	slackArchiveArgChannel  = "channel"
	slackArchiveArgFromDate = "from-date"
	slackArchiveArgToDate   = "to-date"
	slackArchiveArgFilter   = "filter"

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

			var negativeFilters []slackclient.ArchiveFilter
			negatives := map[string]slackclient.ArchiveFilter{"no-link": slackclient.NoLinkFilter{}}
			for _, f := range ctx.StringSlice(slackArchiveArgFilter) {
				filter, ok := negatives[f]
				if !ok {
					return fmt.Errorf("no adequate filter for %s", f)
				}
				negativeFilters = append(negativeFilters, filter)
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
