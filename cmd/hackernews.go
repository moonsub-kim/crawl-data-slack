package main

import (
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/hackernews"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func CrawlHackerNews(c *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger := zapLogger()

	db, err := openMysql(mysqlConn)
	if err != nil {
		return err
	}

	logger.Info("hackernews slack channel", zap.Any("channel", c.String("channel")))
	repository := repository.NewRepository(logger, db)
	hackerNewsCrawler := hackernews.NewCrawler(logger, c.String("channel"), c.Int("point_threshold"))
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		hackerNewsCrawler,
		client,
		client,
	)

	err = usecase.Work(hackerNewsCrawler.GetCrawlerName(), hackerNewsCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
