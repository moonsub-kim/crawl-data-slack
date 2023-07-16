package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/githubclient"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/slackclient"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dateLayout = "2006-01-02"

	envIsDeubg       = "IS_DEBUG"
	envMysqlConn     = "MYSQL_CONN"
	envPostgresConn  = "POSTGRES_CONN"
	envSlackBotToken = "SLACK_BOT_TOKEN"
	envGithubToken   = "GITHUB_TOKEN"

	crawlArgChannel    = "channel"
	crawlArgRecentDays = "recent-days"

	githubArgOwner = "owner"
	githubArgRepo  = "repo"

	Commands = []*cli.Command{
		{
			Name: "maintenance",
			Subcommands: []*cli.Command{
				commandRemoveOldEvents,
				commandMigrateDB,
			},
		},
		{
			Name: "crawl",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: crawlArgChannel},
				&cli.IntFlag{Name: crawlArgRecentDays, DefaultText: "5"},
			},
			Subcommands: []*cli.Command{
				{
					Name: "finance",
					Subcommands: []*cli.Command{
						commandGlobalMonitor,
						commandHankyung,
						commandIPO,
						commandMiraeAsset,
						commanKCIF,
					},
				},
				{
					Name: "tech",
					Subcommands: []*cli.Command{
						commandGoldmanSachs,
						commandHackerNews,
						commandQuastor,
						commandDeliveryHero,
						commandNaverD2,
					},
				},
				{
					Name: "career",
					Subcommands: []*cli.Command{
						commandWanted,
						commandDesignerJob,
						commandNaverCareer,
					},
				},
				commandRSS,
				commandConfluent,
				commandLotteCinema,
				commandInterpark,
			},
		},
		{
			Name: "slack",
			Subcommands: []*cli.Command{
				commandListConversations,
			},
		},
		{
			Name: "github",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: githubArgOwner},
				&cli.StringFlag{Name: githubArgRepo},
			},
			Subcommands: []*cli.Command{
				commandGithubCreateIssue,
				commandArchive,
				commandSyncLabel,
			},
		},
	}
)

func getChromeURL(logger *zap.Logger, chromeHost string) (string, error) {
	endpoint := fmt.Sprintf("http://%s/json/version", chromeHost)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Host = "localhost"

	// request to chrome
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		logger.Error("get", zap.Error(err))
		return "", err
	}
	defer res.Body.Close()

	// read buffer
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("read", zap.Error(err))
		return "", err
	}

	var m map[string]string
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}

	wsURL, ok := m["webSocketDebuggerUrl"]
	if !ok {
		return "", errors.New("webSocketDebuggerUrl is not found")
	}

	u, err := url.Parse(wsURL)
	if err != nil {
		return "", err
	}
	u.Host = chromeHost // replace to chrome host
	return u.String(), nil
}

func zapLogger(ctx *cli.Context) *zap.Logger {
	isDebug := os.Getenv(envIsDeubg) != ""
	// Create logger configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if isDebug {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	// Create logger with configurations
	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		level,
	))

	kv := map[string]string{}
	for _, k := range ctx.FlagNames() {
		kv[k] = ctx.String(k)
	}
	zapLogger.Info(
		"flags",
		zap.Any("flags", kv),
	)

	return zapLogger
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&repository.Event{},
		&repository.Channel{},
	)
}

func openPostgres(conn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openMysql(conn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openDB(logger *zap.Logger) (*gorm.DB, error) {
	postgresConn := os.Getenv(envPostgresConn)
	mysqlConn := os.Getenv(envMysqlConn)

	var f func(string) (*gorm.DB, error)
	var c string
	if postgresConn != "" {
		f = openPostgres
		c = postgresConn
	} else if mysqlConn != "" {
		f = openMysql
		c = mysqlConn
	} else {
		return nil, errors.New("no connection found")
	}

	return f(c)
}

type runGithubCommandFunc func(ctx *cli.Context, logger *zap.Logger, client *githubclient.Client) error

func RunGithub(f runGithubCommandFunc) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		githubToken := os.Getenv(envGithubToken)

		logger := zapLogger(ctx)

		client := githubclient.NewClient(
			logger,
			github.NewClient(
				oauth2.NewClient(
					context.Background(),
					oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: githubToken},
					),
				),
			),
			ctx.String("owner"),
			ctx.String("repo"),
		)

		return f(ctx, logger, client)
	}
}

type initCrawlerFunc func(ctx *cli.Context, logger *zap.Logger, channel string) (crawler.Crawler, error)

func RunCrawl(initCrawler initCrawlerFunc) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		slackBotToken := os.Getenv(envSlackBotToken)
		channel := ctx.String(crawlArgChannel)
		recentDays := ctx.Int(crawlArgRecentDays)
		after := time.Now().Add(time.Duration(-recentDays) * time.Hour * 24)

		logger := zapLogger(ctx)

		c, err := initCrawler(ctx, logger, channel)
		if err != nil {
			return err
		}

		db, err := openDB(logger)
		if err != nil {
			return err
		}

		u := crawler.NewUseCase(
			logger,
			repository.NewRepository(logger, db),
			c,
			slackclient.NewClient(
				logger,
				slack.New(slackBotToken),
				slackBotToken,
				nil,
			),
			nil,
		)

		err = u.Work(after)
		if err != nil {
			logger.Error(
				"Work Error",
				zap.Error(err),
				zap.String("type", reflect.TypeOf(err).String()),
			)
			return err
		}

		logger.Info("Succeeded")
		return nil
	}
}
