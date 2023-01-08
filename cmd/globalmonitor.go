package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/globalmonitor"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandGlobalMonitor *cli.Command = &cli.Command{
		Name: "global-monitor",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return globalmonitor.NewCrawler(
					logger,
					channel,
				)
			},
		),
	}
)
