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

	chgo "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/EpicStep/gdatum"
	clickhouseAdapter "github.com/EpicStep/gdatum/internal/adapters/clickhouse"
	"github.com/EpicStep/gdatum/internal/collector"
	"github.com/EpicStep/gdatum/internal/config"
	"github.com/EpicStep/gdatum/internal/handlers/admin"
	apiHandler "github.com/EpicStep/gdatum/internal/handlers/api"
	clickhouseRepository "github.com/EpicStep/gdatum/internal/infrastructure/repository/clickhouse"
	"github.com/EpicStep/gdatum/internal/infrastructure/server"
	"github.com/EpicStep/gdatum/internal/infrastructure/worker"
	"github.com/EpicStep/gdatum/internal/metrics"
	"github.com/EpicStep/gdatum/internal/utils/buildinfo"
	"github.com/EpicStep/gdatum/internal/utils/migrations"
	"github.com/EpicStep/gdatum/pkg/api"
)

var (
	showVersion   = flag.Bool("version", false, "Print version")
	runMigrations = flag.Bool("migrate", false, "run migrations")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(buildinfo.Get().String())
		return
	}

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

	db, err := openDB(ctx, cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("openDB: %w", err)
	}

	repo := clickhouseAdapter.New(clickhouseRepository.New(db))

	statsHandler := collector.New(repo, metrics.NewCollectorMetrics(prometheus.DefaultRegisterer), logger)
	statsCollectorWorker := worker.New("stats-collector", time.Hour, statsHandler.Handle, logger)

	apiServer, err := api.NewServer(apiHandler.New(repo))
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

func openDB(ctx context.Context, dsn string) (driver.Conn, error) {
	dbOpts, err := chgo.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.ParseDSN: %w", err)
	}

	db, err := chgo.Open(dbOpts)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
