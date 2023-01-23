package main

import (
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/designerjob"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	designerJobArgQuery  string = "query"
	designerJobArgExcept string = "except"

	commandDesignerJob *cli.Command = &cli.Command{
		Name: "designer-job",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: designerJobArgQuery, Required: true},
			&cli.StringSliceFlag{Name: designerJobArgExcept},
		},
		Action: RunCrawl(
			func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error) {
				return designerjob.NewCrawler(
					logger,
					channel,
					ctx.String(designerJobArgQuery),
					ctx.StringSlice(designerJobArgExcept),
				), nil
			},
		),
	}
)
