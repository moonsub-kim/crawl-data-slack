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

func CrawlSpinnaker(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger := zapLogger()

	db, err := openMysql(mysqlConn)
	if err != nil {
		return err
	}

	logger.Info("slack channel", zap.Any("channel", ctx.String("channel")))
	repository := repository.NewRepository(logger, db)
	spinnakerCrawler, err := spinnaker.NewCrawler(logger, ctx.String("channel"), ctx.String("host"), ctx.String("token"))
	if err != nil {
		return err
	}
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
		spinnakerCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(spinnakerCrawler.GetCrawlerName(), spinnakerCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
