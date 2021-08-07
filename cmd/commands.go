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
		Name: "cron",
		Subcommands: []*cli.Command{
			{
				Name: "exclude",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "expression"},
				},
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
		},
	},
}

// CrawlGroupWareDeclinedPayments crawls declied payments from groupware and notify the events
func CrawlGroupWareDeclinedPayments(ctx *cli.Context) error {
	return nil
}

// ExcludeCronExpression excludes
func ExcludeCronExpression(ctx *cli.Context) error {
	return nil
}
