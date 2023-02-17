package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/interpark"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	interparkArgDate string = "date"

	commandInterpark *cli.Command = &cli.Command{
		Name: "interpark",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: interparkArgDate, Required: true},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return interpark.NewCrawler(
					logger,
					channel,
					ctx.String(interparkArgDate),
				), nil
			},
		),
	}
)
