package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/goldmansachs"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	goldmanSachsArgRecentDays string = "recent-days"

	commandGoldmanSachs *cli.Command = &cli.Command{
		Name: "goldman-sachs",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: goldmanSachsArgRecentDays, Required: false},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return goldmansachs.NewCrawler(
					logger,
					channel,
					ctx.Int64(goldmanSachsArgRecentDays),
				), nil
			},
		),
	}
)
