package main

import (
	"os"
	"reflect"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/rss"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	rssFlagNameChannel          string = "channel"
	rssFlagNameName             string = "name"
	rssFlagNameSite             string = "site"
	rssFlagNameCategoryContains string = "category-contains"
	rssFlagNameURLContains      string = "url-contains"

	commandRSS *cli.Command = &cli.Command{
		Name: "rss",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: rssFlagNameChannel, Required: true},
			&cli.StringFlag{Name: rssFlagNameName, Required: true},
			&cli.StringFlag{Name: rssFlagNameSite, Required: true},
			&cli.StringFlag{Name: rssFlagNameCategoryContains},
			&cli.StringFlag{Name: rssFlagNameURLContains},
		},
		Action: crawlRSS,
	}
)

func crawlRSS(ctx *cli.Context) error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	postgresConn := os.Getenv("POSTGRES_CONN")

	logger := zapLogger()

	db, err := openPostgres(postgresConn)
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
	)

	rssCrawler := rss.NewCrawler(
		logger,
		ctx.String(rssFlagNameChannel),
		ctx.String(rssFlagNameName),
		ctx.String(rssFlagNameSite),
		rss.WithCategoryMustContainsFilter(strings.Split(ctx.String(rssFlagNameCategoryContains), ",")),
		rss.WithURLMustContainsFilter(strings.Split(ctx.String(rssFlagNameURLContains), ",")),
	)

	// events, err := rssCrawler.Crawl()
	// logger.Info("end", zap.Any("events", events), zap.Error(err))

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
