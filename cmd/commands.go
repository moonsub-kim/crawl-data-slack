package main

import (
	"os"

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
					&cli.BoolFlag{Name: "job"},
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
	// groupWareID := os.Getenv("GROUPWARE_ID")
	// groupWarePW := os.Getenv("GROUPWARE_PW")
	// slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	log := zaplogger.NewZapLoggerWrapper(zapLogger)

	db, err := gorm.Open(mysql.Open(mysqlConn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&repository.Event{},
		&repository.Restriction{},
		&repository.User{},
	)
	if err != nil {
		return err
	}

	repository := repository.NewRepository(log, db)
	groupwareCrawler := groupware.NewCrawler()
	slackService := slack.NewService()

	usecase := crawler.NewUseCase(
		log,
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
