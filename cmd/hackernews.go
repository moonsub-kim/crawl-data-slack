package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/hackernews"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	hackerNewsArgPointThreshold string = "point_threshold"

	commandHackerNews *cli.Command = &cli.Command{
		Name: "hacker-news",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: hackerNewsArgPointThreshold, Required: false},
		},
		Action: Run(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return hackernews.NewCrawler(
					logger,
					channel,
					ctx.Int(hackerNewsArgPointThreshold),
				), nil
			},
		),
	}
)
