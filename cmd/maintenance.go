package main

import (
	"os"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler/repository"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	removeOldEventsArgTTLDays string = "ttl-days"

	commandRemoveOldEvents *cli.Command = &cli.Command{
		Name: "remove-old-events",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: removeOldEventsArgTTLDays},
		},
		Action: func(ctx *cli.Context) error {
			logger := zapLogger(ctx)

			db, err := openDB(logger)
			if err != nil {
				return err
			}

			ttlDays := ctx.Int(removeOldEventsArgTTLDays)
			before := time.Now().Add(-1 * time.Duration(ttlDays) * 24 * time.Hour)

			r := repository.NewRepository(logger, db)
			cnt, err := r.RemoveOldEvents(before)
			if err != nil {
				return err
			}

			logger.Info(
				"remove",
				zap.Int("count", cnt),
			)

			return nil
		},
	}

	commandMigrateDB *cli.Command = &cli.Command{
		Name:        "migrate-db",
		Description: "It moves all data from old db to new db",
		Action: func(ctx *cli.Context) error {
			oldDB, err := openPostgres(os.Getenv(envPostgresConn))
			if err != nil {
				return err
			}

			newDB, err := openPostgres(os.Getenv("NEW_POSTGRES_CONN"))
			if err != nil {
				return err
			}

			// Create Schema
			migrate(newDB)

			var events []repository.Event
			if oldDB.Find(&events).Error != nil {
				return err
			}

			return newDB.Create(&events).Error
		},
	}
)
