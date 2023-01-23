package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/wanted"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	wantedArgQuery  string = "query"
	wantedArgExcept string = "except"

	commandWanted *cli.Command = &cli.Command{
		Name: "wanted",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: wantedArgQuery, Required: true},
			&cli.StringSliceFlag{Name: wantedArgExcept},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return wanted.NewCrawler(
					logger,
					channel,
					ctx.String(wantedArgQuery),
					ctx.StringSlice(wantedArgExcept),
				), nil
			},
		),
	}
)
