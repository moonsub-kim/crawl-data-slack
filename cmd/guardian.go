package main

import (
	"context"
	"flag"
	"log"
	"os"
	"reflect"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/guardian"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func CrawlGuardian(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")
	chromeHost := os.Getenv("CHROME_HOST")

	logger := zapLogger()

	logger.Info("conn", zap.Any("con", postgresConn))
	db, err := openPostgres(postgresConn)
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
		chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	repository := repository.NewRepository(logger, db)
	guardianCrawler := guardian.NewCrawler(logger, chromectx, ctx.String("channel"), ctx.String("url"))
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
		guardianCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(guardianCrawler.GetCrawlerName(), guardianCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
