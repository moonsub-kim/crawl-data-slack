package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/navercareer"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	naverCareerArgQuery   string = "query"
	naverCareerArgInclude string = "include"
	naverCareerArgExclude string = "exclude"

	commandNaverCareer *cli.Command = &cli.Command{
		Name: "naver-career",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: naverCareerArgQuery, Required: true},
			&cli.StringSliceFlag{Name: naverCareerArgExclude},
			&cli.StringSliceFlag{Name: naverCareerArgInclude},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return navercareer.NewCrawler(
					logger,
					channel,
					ctx.String(naverCareerArgQuery),
					ctx.StringSlice(naverCareerArgInclude),
					ctx.StringSlice(naverCareerArgExclude),
				), nil
			},
		),
	}
)
