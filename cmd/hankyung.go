package main

import (
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/hankyung"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	hankyungFlagNameChannel string = "channel"

	commandHankyung *cli.Command = &cli.Command{
		Name: "hankyung",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: hankyungFlagNameChannel, Required: true},
		},
		Action: CrawlHankyung,
	}
)

func CrawlHankyung(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")

	logger := zapLogger()

	db, err := openPostgres(postgresConn)
	if err != nil {
		return err
	}

	repository := repository.NewRepository(logger, db)
	financialReportCrawler := hankyung.NewCrawler(logger, ctx.String("channel"))
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
		financialReportCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(financialReportCrawler.GetCrawlerName(), financialReportCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}