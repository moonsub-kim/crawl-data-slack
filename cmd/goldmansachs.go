package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/goldmansachs"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandGoldmanSachs *cli.Command = &cli.Command{
		Name: "goldman-sachs",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return goldmansachs.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
