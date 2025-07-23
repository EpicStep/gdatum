// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	ragempAdapter "github.com/EpicStep/gdatum/internal/adapters/ragemp"
	"github.com/EpicStep/gdatum/internal/domain"
	ragempClient "github.com/EpicStep/gdatum/internal/infrastructure/clients/ragemp"
	backoffUtils "github.com/EpicStep/gdatum/internal/utils/backoff"
)

// Handler ...
type Handler struct {
	collectors []collectInstance
	repo       domain.Repository

	collectedGauge       *prometheus.GaugeVec
	collectFailedCounter *prometheus.CounterVec
	insertFailedCounter  prometheus.Counter

	logger *zap.Logger
}

// New ...
func New(repo domain.Repository, logger *zap.Logger) *Handler {
	if logger == nil {
		logger = zap.L()
	}

	ragemp := ragempAdapter.New(ragempClient.New(ragempClient.NewOpts{})) // TODO: make general client to egress

	return &Handler{
		collectors: []collectInstance{
			{
				Multiplayer: domain.MultiplayerRagemp,
				Collect:     ragemp.Servers,
			},
		},
		repo: repo,

		collectedGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "stats_servers_collected_count",
				Help: "Number of servers that have been collected",
			},
			[]string{"multiplayer"}),
		collectFailedCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stats_collect_failed_total",
				Help: "The total number of failed collects of server stats",
			},
			[]string{"multiplayer"}),
		insertFailedCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "stats_insert_failed_total",
				Help: "The total number of failed inserts to repository of server stats",
			},
		),

		logger: logger,
	}
}

// Handle ...
func (h *Handler) Handle(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	servers := h.collect(ctx)

	var insertAttempt int
	_, err := backoff.Retry(
		ctx,
		backoffUtils.EmptyReturnOperation(func() error {
			err := h.repo.InsertServers(ctx, servers)
			if err != nil {
				insertAttempt++
				h.logger.Error("failed to insert servers",
					zap.Int("attempt", insertAttempt),
					zap.Error(err),
				)

				return fmt.Errorf("h.repo.InsertServers: %w", err)
			}

			return nil
		}),
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(3),
	)
	if err != nil {
		h.insertFailedCounter.Inc()
		return fmt.Errorf("backoff.Retry: %w", err)
	}

	return nil
}
