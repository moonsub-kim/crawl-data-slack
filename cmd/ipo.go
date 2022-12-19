package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/ipo"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandIPO *cli.Command = &cli.Command{
		Name: "ipo",
		Action: Run(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return ipo.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
