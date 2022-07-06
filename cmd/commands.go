package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Commands = []*cli.Command{
	{
		Name: "crawl",
		Subcommands: []*cli.Command{
			{
				Name: "groupware",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "job"},
					&cli.StringFlag{Name: "masters"},
					&cli.StringFlag{Name: "renames"},
				},
				Action: CrawlGroupWareDeclinedPayments,
			},
			{
				Name: "hackernews",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.IntFlag{Name: "point_threshold"},
				},
				Action: CrawlHackerNews,
			},
			{
				Name: "quasarzone",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlQuasarZoneSales,
			},
			{
				Name: "gitpublic",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "organization"},
				},
				Action: CrawlGitPublic,
			},
			{
				Name: "wanted",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "query"},
				},
				Action: CrawlWanted,
			},
			{
				Name: "eomisae",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "target"},
				},
				Action: CrawlEomisae,
			},
			{
				Name: "ipo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlIPO,
			},
			{
				Name: "financial-report",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlFinancialReport,
			},
			{
				Name: "spinnaker",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "host"},
					&cli.StringFlag{Name: "token"},
				},
				Action: CrawlSpinnaker,
			},
			{
				Name: "techcrunch",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlTechCrunch,
			},
			{
				Name: "confluent",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlConfluent,
			},
			commandRSS,
			commandSlackEngineering,
			commandShopifyEngineering,
			commandMiraeAsset,
		},
	},
	{
		Name: "slack",
		Subcommands: []*cli.Command{
			{
				Name: "get-channel",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: GetChannel,
			},
		},
	},
}

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
	body, err := ioutil.ReadAll(res.Body)
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

func zapLogger() *zap.Logger {
	// Create logger configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create logger with configurations
	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	))

	return zapLogger
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&repository.Event{},
		&repository.Channel{},
	)
}

func toRenameMap(l *zap.Logger, s string) (map[string]string, error) {
	if s == "" {
		return map[string]string{}, nil
	}

	m := map[string]string{}
	renames := strings.Split(s, ",")
	for _, rename := range renames {
		splitted := strings.Split(s, "=")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("failed to parse rename %s", rename)
		}

		m[splitted[0]] = splitted[1]
	}
	l.Info("renameMap", zap.Any("map", m))
	return m, nil
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
