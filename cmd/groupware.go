package main

import (
	"context"
	"flag"
	"os"
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/groupwaredecline"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// CrawlGroupWareDeclinedPayments crawls declied payments from groupware and notify the events
func CrawlGroupWareDeclinedPayments(ctx *cli.Context) error {
	groupWareID := os.Getenv("GROUPWARE_ID")
	groupWarePW := os.Getenv("GROUPWARE_PW")
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")
	chromeHost := os.Getenv("CHROME_HOST")

	logger := zapLogger()

	m, err := toRenameMap(logger, ctx.String("renames"))
	if err != nil {
		logger.Error("", zap.Error(err))
		return err
	}

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

	var masters []string
	if ctx.String("masters") != "" {
		masters = strings.Split(ctx.String("masters"), ",")
		logger.Info("masters", zap.Any("masters", masters))
	}

	repository := repository.NewRepository(logger, db)
	groupwareCrawler := groupwaredecline.NewCrawler(logger, chromectx, groupWareID, groupWarePW, masters)
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		groupwareCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(groupwareCrawler.GetCrawlerName(), groupwareCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
