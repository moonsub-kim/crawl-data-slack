package main

import (
	"context"
	"flag"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/confluent"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	confluentArgJob     string = "job"
	confluentArgKeyword string = "keyword"

	commandConfluent *cli.Command = &cli.Command{
		Name: "confluent",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: confluentArgJob, Required: true},
			&cli.StringSliceFlag{Name: confluentArgKeyword, Required: false, Usage: "space separated keywords"},
		},
		Action: Run(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				chromeHost := os.Getenv("CHROME_HOST")
				url, err := getChromeURL(logger, chromeHost)
				if err != nil {
					return nil, err
				}
				logger.Info("chrome url", zap.String("url", url))

				devtoolsWSURL := flag.String("devtools-ws-url", url, "DevTools Websocket URL")
				allocatorctx, cancel := chromedp.NewRemoteAllocator(context.Background(), *devtoolsWSURL)
				defer cancel()

				chromectx, cancel := chromedp.NewContext(
					allocatorctx,
					// chromedp.WithLogf(log.Printf),
					// chromedp.WithDebugf(log.Printf),
				)
				defer cancel()

				return confluent.NewCrawler(
					logger,
					chromectx,
					channel,
					ctx.String(confluentArgJob),
					ctx.StringSlice(confluentArgKeyword),
				), nil
			},
		),
	}
)
