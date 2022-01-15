package main

import (
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/quasarzone"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func CrawlQuasarZoneSales(c *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")

	logger := zapLogger()

	db, err := openPostgres(postgresConn)
	if err != nil {
		return err
	}

	logger.Info("slack channel", zap.Any("channel", c.String("channel")))
	repository := repository.NewRepository(logger, db)
	quasarZoneCrawler := quasarzone.NewCrawler(logger, c.String("channel"))
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		quasarZoneCrawler,
		client,
		client,
	)

	err = usecase.Work(quasarZoneCrawler.GetCrawlerName(), quasarZoneCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
