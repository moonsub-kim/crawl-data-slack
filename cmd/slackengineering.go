package main

import (
	"errors"
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackengineering"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	slackEngineeringFlagNameChannel string = "channel"

	commandSlackEngineering *cli.Command = &cli.Command{
		Name: "slack-engineering",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: slackEngineeringFlagNameChannel, Required: true},
		},
		Action: crawlSlackEngineering,
	}
)

func crawlSlackEngineering(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger := zapLogger()

	var f func(string) (*gorm.DB, error)
	var c string
	if postgresConn != "" {
		f = openPostgres
		c = postgresConn
	} else if mysqlConn != "" {
		f = openMysql
		c = mysqlConn
	} else {
		return errors.New("no connection found")
	}

	db, err := f(c)
	if err != nil {
		return err
	}

	logger.Info(
		"args",
		zap.Any(slackEngineeringFlagNameChannel, ctx.String(slackEngineeringFlagNameChannel)),
	)

	rssCrawler := slackengineering.NewCrawler(
		logger,
		ctx.String(rssFlagNameChannel),
	)

	events, err := rssCrawler.Crawl()
	logger.Info("result", zap.Any("events", events), zap.Error(err))

	repository := repository.NewRepository(logger, db)
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)
	m, err := toRenameMap(logger, ctx.String("renames"))
	if err != nil {
		logger.Error("", zap.Error(err))
		return err
	}

	usecase := crawler.NewUseCase(
		logger,
		repository,
		rssCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(rssCrawler.GetCrawlerName(), rssCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
