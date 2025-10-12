// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"

	"github.com/EpicStep/gdatum/internal/domain"
	"github.com/EpicStep/gdatum/internal/infrastructure/repository/clickhouse"
)

type clickhouseStore interface {
	InsertServers(ctx context.Context, servers []clickhouse.Server) error
	ListMultiplayerSummaries(ctx context.Context, playersOrderAsc bool) ([]clickhouse.MultiplayerSummary, error)
	ListServerSummaries(ctx context.Context, params domain.ListServerSummariesParams) ([]clickhouse.ServerSummary, error)
	GetServer(ctx context.Context, multiplayer domain.Multiplayer, host string) (clickhouse.Server, error)
	ListServerStatistics(ctx context.Context, params domain.ListServerStatisticsParams) ([]clickhouse.ServerStatisticPoint, error)
}

// Adapter ...
type Adapter struct {
	store clickhouseStore
}

// New returns new ClickHouse adapter.
func New(store clickhouseStore) *Adapter {
	return &Adapter{
		store: store,
	}
}

// InsertServers ...
func (a *Adapter) InsertServers(ctx context.Context, servers []domain.Server) error {
	chServers := lo.Map(servers, func(srv domain.Server, _ int) clickhouse.Server {
		return clickhouse.Server{
			Multiplayer:  string(srv.Multiplayer),
			Host:         srv.Host,
			Name:         srv.Name,
			URL:          srv.URL,
			Gamemode:     srv.Gamemode,
			Language:     srv.Language,
			PlayersCount: srv.PlayersCount,
			CollectedAt:  srv.CollectedAt,
		}
	})

	if err := a.store.InsertServers(ctx, chServers); err != nil {
		return err
	}

	return nil
}

// ListMultiplayerSummaries ...
func (a *Adapter) ListMultiplayerSummaries(ctx context.Context, playersOrderAsc bool) ([]domain.MultiplayerSummary, error) {
	summaries, err := a.store.ListMultiplayerSummaries(ctx, playersOrderAsc)
	if err != nil {
		return nil, err
	}

	return lo.Map(summaries, func(summary clickhouse.MultiplayerSummary, _ int) domain.MultiplayerSummary {
		return domain.MultiplayerSummary{
			Name:         domain.Multiplayer(summary.Multiplayer),
			PlayersCount: summary.PlayersCount,
		}
	}), nil
}

// GetServer ...
func (a *Adapter) GetServer(ctx context.Context, multiplayer domain.Multiplayer, host string) (domain.Server, error) {
	chServer, err := a.store.GetServer(ctx, multiplayer, host)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Server{}, domain.ErrServerNotFound
		}

		return domain.Server{}, err
	}

	return domain.Server{
		Multiplayer:  domain.Multiplayer(chServer.Multiplayer),
		Host:         chServer.Host,
		Name:         chServer.Name,
		URL:          chServer.URL,
		Gamemode:     chServer.Gamemode,
		Language:     chServer.Language,
		PlayersCount: chServer.PlayersCount,
		CollectedAt:  chServer.CollectedAt,
	}, nil
}

// ListServerSummaries ...
func (a *Adapter) ListServerSummaries(ctx context.Context, params domain.ListServerSummariesParams) ([]domain.ServerSummary, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("params.Validate: %w", err)
	}

	servers, err := a.store.ListServerSummaries(ctx, params)
	if err != nil {
		return nil, err
	}

	return lo.Map(servers, func(server clickhouse.ServerSummary, _ int) domain.ServerSummary {
		return domain.ServerSummary{
			Host:         server.Host,
			Name:         server.Name,
			PlayersCount: server.PlayersCount,
		}
	}), nil
}

// ListServerStatistics ...
func (a *Adapter) ListServerStatistics(ctx context.Context, params domain.ListServerStatisticsParams) ([]domain.ServerStatisticPoint, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("params.Validate: %w", err)
	}

	statistics, err := a.store.ListServerStatistics(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(statistics) == 0 {
		return nil, domain.ErrServerNotFound
	}

	return lo.Map(statistics, func(statistic clickhouse.ServerStatisticPoint, _ int) domain.ServerStatisticPoint {
		return domain.ServerStatisticPoint{
			PlayersCount: statistic.PlayersCount,
			CollectedAt:  statistic.CollectedAt,
		}
	}), nil
}
