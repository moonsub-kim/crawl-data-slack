package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/quastor"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	quastorArgRecentDays string = "recent-days"

	commandQuastor *cli.Command = &cli.Command{
		Name: "quastor",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return quastor.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
