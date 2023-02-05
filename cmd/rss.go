package main

import (
	"strings"

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
	rssArgFetchRSS         string = "fetch-rss"
	rssArgTechBlogPosts    string = "tech-blog-posts"

	commandRSS *cli.Command = &cli.Command{
		Name: "rss",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: rssArgName, Required: true},
			&cli.StringFlag{Name: rssArgSite, Required: true},
			&cli.StringFlag{Name: rssArgCategoryContains},
			&cli.StringFlag{Name: rssArgURLContains},
			&cli.BoolFlag{Name: rssArgFetchRSS},
			&cli.BoolFlag{Name: rssArgTechBlogPosts},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				var opts []rss.CrawlerOption
				if urlContains := ctx.String(rssArgURLContains); urlContains != "" {
					opts = append(opts, rss.WithURLMustContainsTransformer(strings.Split(urlContains, ",")))
				}

				if categoryContains := ctx.String(rssArgCategoryContains); categoryContains != "" {
					opts = append(opts, rss.WithCategoryMustContainsTransformer(strings.Split(categoryContains, ",")))
				}

				if ctx.Bool(rssArgFetchRSS) {
					opts = append(opts, rss.WithFetchRSSTransformer())
				}

				if ctx.Bool(rssArgTechBlogPosts) {
					opts = append(opts, rss.WithTechBlogPostsTransformer())
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
