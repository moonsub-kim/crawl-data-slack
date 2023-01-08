package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/deliveryhero"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	commandDeliveryHero *cli.Command = &cli.Command{
		Name: "delivery-hero",
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return deliveryhero.NewCrawler(
					logger,
					channel,
				), nil
			},
		),
	}
)
