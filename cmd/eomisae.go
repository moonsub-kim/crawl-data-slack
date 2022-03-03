package main

import (
	"context"
	"flag"
	"os"
	"reflect"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/eomisae"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// CrawlEomisae crawls declied payments from groupware and notify the events
func CrawlEomisae(ctx *cli.Context) error {
	id := os.Getenv("EOMISAE_ID")
	pw := os.Getenv("EOMISAE_PW")
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")
	chromeHost := os.Getenv("CHROME_HOST")
	channel := ctx.String("channel")
	target := ctx.String("target")

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

	repository := repository.NewRepository(logger, db)
	eomisaeCrawler, err := eomisae.NewCrawler(logger, chromectx, channel, target, id, pw)
	if err != nil {
		logger.Error("", zap.Error(err))
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
		eomisaeCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(eomisaeCrawler.GetCrawlerName(), eomisaeCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
