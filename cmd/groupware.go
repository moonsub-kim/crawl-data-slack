package main

import (
	"context"
	"flag"
	"os"
	"reflect"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/groupwaredecline"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CrawlGroupWareDeclinedPayments crawls declied payments from groupware and notify the events
func CrawlGroupWareDeclinedPayments(ctx *cli.Context) error {
	groupWareID := os.Getenv("GROUPWARE_ID")
	groupWarePW := os.Getenv("GROUPWARE_PW")
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	mysqlConn := os.Getenv("MYSQL_CONN")
	chromeHost := os.Getenv("CHROME_HOST")

	logger := zapLogger()

	db, err := gorm.Open(mysql.Open(mysqlConn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&repository.Event{},
		&repository.Restriction{},
		&repository.User{},
	)
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
	groupwareCrawler := groupwaredecline.NewCrawler(logger, chromectx, groupWareID, groupWarePW)
	api := slack.New(slackBotToken)
	client := slackclient.NewClient(logger, api)

	usecase := crawler.NewUseCase(
		logger,
		repository,
		groupwareCrawler,
		client,
		client,
	)

	err = usecase.Work(groupwareCrawler.GetCrawlerName(), groupwareCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}