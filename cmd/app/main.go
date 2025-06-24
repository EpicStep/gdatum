// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/EpicStep/gdatum"
	"github.com/EpicStep/gdatum/internal/config"
	"github.com/EpicStep/gdatum/internal/handlers/admin"
	apiHandler "github.com/EpicStep/gdatum/internal/handlers/api"
	"github.com/EpicStep/gdatum/internal/repository"
	"github.com/EpicStep/gdatum/internal/stats"
	"github.com/EpicStep/gdatum/internal/utils/migrations"
	"github.com/EpicStep/gdatum/internal/utils/server"
	"github.com/EpicStep/gdatum/internal/worker"
	"github.com/EpicStep/gdatum/pkg/api"
)

var (
	runMigrations = flag.Bool("migrate", false, "run migrations")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction() //nolint:errcheck
	defer logger.Sync()              //nolint:errcheck
	zap.ReplaceGlobals(logger)

	if err := run(logger); err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}

func run(logger *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config.Load: %w", err)
	}

	zap.L().Info("loaded config", zap.Inline(cfg))

	if *runMigrations {
		logger.Info("running migrations")
		err = migrations.Run(ctx, cfg.DatabaseDSN, gdatum.MigrationsFS)
		if err != nil {
			return fmt.Errorf("migrations.Run: %w", err)
		}

		return nil
	}

	db, err := openDB(cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("openDB: %w", err)
	}

	repo := repository.New(db)

	statsHandler := stats.New(repo, logger)
	statsCollectorWorker := worker.New("stats-collector", time.Hour, statsHandler.Handle, logger)

	apiServer, err := api.NewServer(apiHandler.New())
	if err != nil {
		return fmt.Errorf("api.NewServer: %w", err)
	}

	publicServer := server.New(cfg.PublicListenAddress, apiServer, logger.With(zap.String("kind", "public")))
	adminServer := server.New(cfg.AdminListenAddress, admin.Handler(), logger.With(zap.String("kind", "admin")))

	eg, eCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return publicServer.Run(eCtx)
	})
	eg.Go(func() error {
		return adminServer.Run(eCtx)
	})
	eg.Go(func() error {
		return statsCollectorWorker.Run(eCtx)
	})

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("eg.Wait: %w", err)
	}

	return nil
}

func openDB(dsn string) (driver.Conn, error) {
	dbOpts, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.ParseDSN: %w", err)
	}

	db, err := clickhouse.Open(dbOpts)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.Open: %w", err)
	}

	return db, nil
}
