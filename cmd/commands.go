package main

import (
	"os"
	"reflect"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
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
	{
		Name: "test",
		Subcommands: []*cli.Command{
			{Name: "slack", Action: TestSlack},
			{Name: "chrome", Action: TestChrome},
		},
	},
}

// CrawlGroupWareDeclinedPayments crawls declied payments from groupware and notify the events
func CrawlGroupWareDeclinedPayments(ctx *cli.Context) error {
	// groupWareID := os.Getenv("GROUPWARE_ID")
	// groupWarePW := os.Getenv("GROUPWARE_PW")
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	logger.Info("CrawlGroupWareDeclinedPayments")

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

	repository := repository.NewRepository(logger, db)
	var groupwareCrawler crawler.Crawler // := groupware.NewCrawler(log)
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		groupwareCrawler,
		client,
		client,
	)

	err = usecase.Work("groupware", "declined_payments")
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}
	return nil
}

// AddRestriction adds a restriction
func AddRestriction(ctx *cli.Context) error {
	return nil
}

func TestSlack(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	client.Notify(crawler.Notification{
		User: crawler.User{
			ID: "UJBG25A04",
		},
		Event: crawler.Event{
			Crawler: "groupware",
			Job:     "declined_payments",
		},
	})

	return nil
}

func TestChrome(ctx *cli.Context) error {
	return nil
}
