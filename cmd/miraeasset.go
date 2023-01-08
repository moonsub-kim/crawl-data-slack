package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/miraeasset"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandMiraeAsset *cli.Command = &cli.Command{
		Name: "mirae-asset",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return miraeasset.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
