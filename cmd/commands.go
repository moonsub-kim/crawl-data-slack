package main

import (
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Name: "crawl",
		Subcommands: []*cli.Command{
			{
				Name: "groupware",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "declined_payments"},
				},
				Action: CrawlGroupWareDeclinedPayments,
			},
		},
	},
	{
		Name: "restriction",
		Subcommands: []*cli.Command{
			{
				Name: "add",
				Flags: []cli.Flag{
					&cli.TimestampFlag{Name: "start_date"},
					&cli.TimestampFlag{Name: "end_date"},
					&cli.TimestampFlag{Name: "hour_from"},
					&cli.TimestampFlag{Name: "hour_to"},
				},
				Action: AddRestriction,
			},
		},
	},
}

// CrawlGroupWareDeclinedPayments crawls declied payments from groupware and notify the events
func CrawlGroupWareDeclinedPayments(ctx *cli.Context) error {
	return nil
}

// AddRestriction adds a restriction
func AddRestriction(ctx *cli.Context) error {
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	return err
	// }
	// var usecase Usecase = crawler.NewUseCase(
	// 	logger,
	// )
	return nil
}
