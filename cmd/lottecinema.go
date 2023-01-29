package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/lottecinema"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	lotteCinemaArgDate string = "date"

	commandLotteCinema *cli.Command = &cli.Command{
		Name: "lotte-cinema",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: lotteCinemaArgDate, Required: true},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return lottecinema.NewCrawler(
					logger,
					channel,
					ctx.String(lotteCinemaArgDate),
				), nil
			},
		),
	}
)
