// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"github.com/samber/lo"

	"github.com/EpicStep/gdatum/internal/domain"
	"github.com/EpicStep/gdatum/internal/infrastructure/repository/clickhouse"
)

type clickhouseStore interface {
	InsertServers(ctx context.Context, servers []clickhouse.Server) error
	GetMultiplayersSummary(ctx context.Context, playersOrder domain.Order) ([]clickhouse.MultiplayerSummary, error)
	GetServerByID(ctx context.Context, multiplayer domain.Multiplayer, id string) (clickhouse.Server, error)
	GetServersByMultiplayer(ctx context.Context, filter domain.ServersByMultiplayerFilter) ([]clickhouse.ServerSummary, error)
	GetServerStatistics(ctx context.Context, filter domain.ServerStatisticsFilter) ([]clickhouse.ServerStatistic, error)
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
			Multiplayer: string(srv.Multiplayer),
			ID:          srv.ID,
			Name:        srv.Name,
			URL:         srv.URL,
			Gamemode:    srv.Gamemode,
			Lang:        srv.Lang,
			Players:     srv.Players,
			CollectedAt: srv.CollectedAt,
		}
	})

	if err := a.store.InsertServers(ctx, chServers); err != nil {
		return err
	}

	return nil
}

// GetMultiplayersSummary ...
func (a *Adapter) GetMultiplayersSummary(ctx context.Context, playersOrder domain.Order) ([]domain.MultiplayerSummary, error) {
	summaries, err := a.store.GetMultiplayersSummary(ctx, playersOrder)
	if err != nil {
		return nil, err
	}

	return lo.Map(summaries, func(summary clickhouse.MultiplayerSummary, _ int) domain.MultiplayerSummary {
		return domain.MultiplayerSummary{
			Name:    domain.Multiplayer(summary.Multiplayer),
			Players: summary.Players,
		}
	}), nil
}

// GetServerByID ...
func (a *Adapter) GetServerByID(ctx context.Context, multiplayer domain.Multiplayer, id string) (domain.Server, error) {
	chServer, err := a.store.GetServerByID(ctx, multiplayer, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Server{}, domain.ErrServerNotFound
		}

		return domain.Server{}, err
	}

	return domain.Server{
		Multiplayer: domain.Multiplayer(chServer.Multiplayer),
		ID:          chServer.ID,
		Name:        chServer.Name,
		URL:         chServer.URL,
		Gamemode:    chServer.Gamemode,
		Lang:        chServer.Lang,
		Players:     chServer.Players,
		CollectedAt: chServer.CollectedAt,
	}, nil
}

// GetServersByMultiplayer ...
func (a *Adapter) GetServersByMultiplayer(ctx context.Context, filter domain.ServersByMultiplayerFilter) ([]domain.ServerSummary, error) {
	servers, err := a.store.GetServersByMultiplayer(ctx, filter)
	if err != nil {
		return nil, err
	}

	return lo.Map(servers, func(server clickhouse.ServerSummary, _ int) domain.ServerSummary {
		return domain.ServerSummary{
			ID:      server.ID,
			Name:    server.Name,
			Players: server.Players,
		}
	}), nil
}

// GetServerStatistics ...
func (a *Adapter) GetServerStatistics(ctx context.Context, filter domain.ServerStatisticsFilter) ([]domain.ServerStatistic, error) {
	statistics, err := a.store.GetServerStatistics(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(statistics) == 0 {
		return nil, domain.ErrServerNotFound
	}

	return lo.Map(statistics, func(statistic clickhouse.ServerStatistic, _ int) domain.ServerStatistic {
		return domain.ServerStatistic{
			Players: statistic.Players,
			At:      statistic.At,
		}
	}), nil
}
