package main

import (
	"context"
	"flag"
	"os"
	"reflect"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/confluent"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	confluentFlagNameChannel string = "channel"
	confluentFlagNameJob     string = "job"
	confluentFlagNameKeyword string = "keyword"

	commandConfluent *cli.Command = &cli.Command{
		Name: "confluent",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: confluentFlagNameChannel, Required: true},
			&cli.StringFlag{Name: confluentFlagNameJob, Required: true},
			&cli.StringSliceFlag{Name: confluentFlagNameKeyword, Required: false, Usage: "space separated keywords"},
		},
		Action: CrawlConfluent,
	}
)

func CrawlConfluent(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")
	chromeHost := os.Getenv("CHROME_HOST")

	logger := zapLogger()

	db, err := openMysql(mysqlConn)
	if err != nil {
		return err
	}

	url, err := getChromeURL(logger, chromeHost)
	if err != nil {
		return err
	}
	logger.Info("chrome url", zap.String("url", url))

	devtoolsWSURL := flag.String("devtools-ws-url", url, "DevTools Websocket URL")
	allocatorctx, cancel := chromedp.NewRemoteAllocator(context.Background(), *devtoolsWSURL)
	defer cancel()

	chromectx, cancel := chromedp.NewContext(
		allocatorctx,
		// chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	logger.Info(
		"confluent",
		zap.Any(confluentFlagNameChannel, ctx.String(confluentFlagNameChannel)),
		zap.Any(confluentFlagNameJob, ctx.String(confluentFlagNameJob)),
		zap.Any(confluentFlagNameKeyword, ctx.StringSlice(confluentFlagNameKeyword)),
	)

	repository := repository.NewRepository(logger, db)
	confluentCrawler := confluent.NewCrawler(
		logger,
		chromectx,
		ctx.String(confluentFlagNameChannel),
		ctx.String(confluentFlagNameJob),
		ctx.StringSlice(confluentFlagNameKeyword),
	)
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
		confluentCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(confluentCrawler.GetCrawlerName(), confluentCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
