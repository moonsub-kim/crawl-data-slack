package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/naverd2"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandNaverD2 *cli.Command = &cli.Command{
		Name:  "naver-d2",
		Flags: []cli.Flag{},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return naverd2.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
