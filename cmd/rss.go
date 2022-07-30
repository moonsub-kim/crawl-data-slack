package main

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/rss"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	rssFlagNameChannel          string = "channel"
	rssFlagNameName             string = "name"
	rssFlagNameSite             string = "site"
	rssFlagNameCategoryContains string = "category-contains"
	rssFlagNameURLContains      string = "url-contains"
	rssFlagNameRecentDays       string = "recent-days"

	commandRSS *cli.Command = &cli.Command{
		Name: "rss",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: rssFlagNameChannel, Required: true},
			&cli.StringFlag{Name: rssFlagNameName, Required: true},
			&cli.StringFlag{Name: rssFlagNameSite, Required: true},
			&cli.StringFlag{Name: rssFlagNameCategoryContains},
			&cli.StringFlag{Name: rssFlagNameURLContains},
			&cli.Int64Flag{Name: rssFlagNameRecentDays},
		},
		Action: crawlRSS,
	}
)

func crawlRSS(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")
	mysqlConn := os.Getenv("MYSQL_CONN")

	logger := zapLogger()

	var f func(string) (*gorm.DB, error)
	var c string
	if postgresConn != "" {
		f = openPostgres
		c = postgresConn
	} else if mysqlConn != "" {
		f = openMysql
		c = mysqlConn
	} else {
		return errors.New("no connection found")
	}

	db, err := f(c)
	if err != nil {
		return err
	}

	logger.Info(
		"args",
		zap.Any(rssFlagNameChannel, ctx.String(rssFlagNameChannel)),
		zap.Any(rssFlagNameName, ctx.String(rssFlagNameName)),
		zap.Any(rssFlagNameSite, ctx.String(rssFlagNameSite)),
		zap.Any(rssFlagNameCategoryContains, ctx.String(rssFlagNameCategoryContains)),
		zap.Any(rssFlagNameURLContains, ctx.String(rssFlagNameURLContains)),
		zap.Any(rssFlagNameRecentDays, ctx.Int64(rssFlagNameRecentDays)),
	)

	var opts []rss.CrawlerOption
	if urlContains := ctx.String(rssFlagNameURLContains); urlContains != "" {
		opts = append(opts, rss.WithURLMustContainsFilter(strings.Split(urlContains, ",")))
	}
	if categoryContains := ctx.String(rssFlagNameCategoryContains); categoryContains != "" {
		opts = append(opts, rss.WithCategoryMustContainsFilter(strings.Split(categoryContains, ",")))
	}
	if recent := ctx.Int64(rssFlagNameRecentDays); recent != 0 {
		t := time.Now().Add(time.Duration(-recent) * time.Hour * 24)
		opts = append(opts, rss.WithRecentFilter(t))
	}

	rssCrawler := rss.NewCrawler(
		logger,
		ctx.String(rssFlagNameChannel),
		ctx.String(rssFlagNameName),
		ctx.String(rssFlagNameSite),
		opts...,
	)

	repository := repository.NewRepository(logger, db)
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
		rssCrawler,
		client,
		client,
		m,
	)

	err = usecase.Work(rssCrawler.GetCrawlerName(), rssCrawler.GetJobName())
	if err != nil {
		logger.Error("Work Error", zap.Error(err), zap.String("type", reflect.TypeOf(err).String()))
		return err
	}

	logger.Info("Succeed")
	return nil
}
