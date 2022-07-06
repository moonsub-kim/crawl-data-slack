package main

import (
	"os"
	"reflect"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/miraeasset"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	miraeAssetFlagNameChannel string = "channel"

	commandMiraeAsset *cli.Command = &cli.Command{
		Name: "mirae-asset",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: miraeAssetFlagNameChannel, Required: true},
		},
		Action: CrawlMiraeAsset,
	}
)

// CrawlMiraeAsset
func CrawlMiraeAsset(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")

	logger := zapLogger()

	db, err := openPostgres(postgresConn)
	if err != nil {
		return err
	}

	repository := repository.NewRepository(logger, db)
	financialReportCrawler := miraeasset.NewCrawler(logger, ctx.String("channel"))
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
