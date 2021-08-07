package main

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/groupware"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/slack"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/zaplogger"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	logger := zaplogger.NewZapLoggerWrapper(zapLogger)

	db, err := gorm.Open(mysql.Open(""), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(
		&repository.Event{},
		&repository.Restriction{},
		&repository.User{},
	)
	repository := repository.NewRepository(logger, db)
	groupwareCrawler := groupware.NewCrawler()
	slackService := slack.NewService()

	usecase := crawler.NewUseCase(
		logger,
		repository,
		groupwareCrawler,
		slackService,
		slackService,
	)

	usecase.Work("groupware", "declined_payments")
	return nil
}

// AddRestriction adds a restriction
func AddRestriction(ctx *cli.Context) error {
	return nil
}
