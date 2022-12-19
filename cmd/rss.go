package main

import (
	"strings"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/rss"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	rssArgName             string = "name"
	rssArgSite             string = "site"
	rssArgCategoryContains string = "category-contains"
	rssArgURLContains      string = "url-contains"
	rssArgRecentDays       string = "recent-days"

	commandRSS *cli.Command = &cli.Command{
		Name: "rss",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: rssArgName, Required: true},
			&cli.StringFlag{Name: rssArgSite, Required: true},
			&cli.StringFlag{Name: rssArgCategoryContains},
			&cli.StringFlag{Name: rssArgURLContains},
			&cli.Int64Flag{Name: rssArgRecentDays},
		},
		Action: Run(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				var opts []rss.CrawlerOption
				if urlContains := ctx.String(rssArgURLContains); urlContains != "" {
					opts = append(opts, rss.WithURLMustContainsFilter(strings.Split(urlContains, ",")))
				}

				if categoryContains := ctx.String(rssArgCategoryContains); categoryContains != "" {
					opts = append(opts, rss.WithCategoryMustContainsFilter(strings.Split(categoryContains, ",")))
				}

				if recent := ctx.Int64(rssArgRecentDays); recent != 0 {
					t := time.Now().Add(time.Duration(-recent) * time.Hour * 24)
					opts = append(opts, rss.WithRecentFilter(t))
				}

				return rss.NewCrawler(
					logger,
					channel,
					ctx.String(rssArgName),
					ctx.String(rssArgSite),
					opts...,
				), nil
			},
		),
	}
)
