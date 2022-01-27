package main

import (
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/spinnaker"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func CrawlSpinnaker(c *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger := zapLogger()

	db, err := openMysql(mysqlConn)
	if err != nil {
		return err
	}

	logger.Info("slack channel", zap.Any("channel", c.String("channel")))
	repository := repository.NewRepository(logger, db)
	spinnakerCrawler, err := spinnaker.NewCrawler(logger, c.String("channel"), c.String("host"), c.String("token"))
	if err != nil {
		return err
	}
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		spinnakerCrawler,
		client,
		client,
	)

	err = usecase.Work(spinnakerCrawler.GetCrawlerName(), spinnakerCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
