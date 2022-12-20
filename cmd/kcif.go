package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/kcif"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commanKCIF *cli.Command = &cli.Command{
		Name: "kcif",
		Action: Run(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return kcif.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
