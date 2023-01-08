package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/hankyung"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandHankyung *cli.Command = &cli.Command{
		Name: "hankyung",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return hankyung.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
