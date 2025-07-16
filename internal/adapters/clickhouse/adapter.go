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
	MultiplayersSummary(ctx context.Context, playersOrder domain.Order) ([]clickhouse.MultiplayerSummary, error)
	ServerByIdentifier(ctx context.Context, multiplayer domain.Multiplayer, identifier string) (clickhouse.Server, error)
	ServersByMultiplayer(ctx context.Context, filter domain.ServersByMultiplayerFilter) ([]clickhouse.ServerSummary, error)
	ServerStatistics(ctx context.Context, filter domain.ServerStatisticsFilter) ([]clickhouse.ServerStatistic, error)
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
			Identifier:  srv.Identifier,
			Name:        srv.Name,
			URL:         srv.URL,
			Gamemode:    srv.Gamemode,
			Lang:        srv.Lang,
			Players:     srv.Players,
			Timestamp:   srv.CollectedAt,
		}
	})

	if err := a.store.InsertServers(ctx, chServers); err != nil {
		return fmt.Errorf("a.store.InsertServers: %w", err)
	}

	return nil
}

// MultiplayersSummary ...
func (a *Adapter) MultiplayersSummary(ctx context.Context, playersOrder domain.Order) ([]domain.MultiplayerSummary, error) {
	summaries, err := a.store.MultiplayersSummary(ctx, playersOrder)
	if err != nil {
		return nil, fmt.Errorf("a.store.GetMultiplayersSummary: %w", err)
	}

	return lo.Map(summaries, func(summary clickhouse.MultiplayerSummary, _ int) domain.MultiplayerSummary {
		return domain.MultiplayerSummary{
			Name:    domain.Multiplayer(summary.Multiplayer),
			Players: summary.Players,
		}
	}), nil
}

// ServerByIdentifier ...
func (a *Adapter) ServerByIdentifier(ctx context.Context, multiplayer domain.Multiplayer, identifier string) (domain.Server, error) {
	chServer, err := a.store.ServerByIdentifier(ctx, multiplayer, identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Server{}, domain.ErrServerNotFound
		}

		return domain.Server{}, fmt.Errorf("a.store.GetServerByID: %w", err)
	}

	return domain.Server{
		Multiplayer: domain.Multiplayer(chServer.Multiplayer),
		Identifier:  chServer.Identifier,
		Name:        chServer.Name,
		URL:         chServer.URL,
		Gamemode:    chServer.Gamemode,
		Lang:        chServer.Lang,
		Players:     chServer.Players,
		CollectedAt: chServer.Timestamp,
	}, nil
}

// ServersByMultiplayer ...
func (a *Adapter) ServersByMultiplayer(ctx context.Context, filter domain.ServersByMultiplayerFilter) ([]domain.ServerSummary, error) {
	servers, err := a.store.ServersByMultiplayer(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("a.store.GetServersByMultiplayer: %w", err)
	}

	return lo.Map(servers, func(server clickhouse.ServerSummary, _ int) domain.ServerSummary {
		return domain.ServerSummary{
			Identifier: server.Identifier,
			Name:       server.Name,
			Players:    server.Players,
		}
	}), nil
}

// ServerStatistics ...
func (a *Adapter) ServerStatistics(ctx context.Context, filter domain.ServerStatisticsFilter) ([]domain.ServerStatistic, error) {
	statistics, err := a.store.ServerStatistics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("a.store.GetServerStatistics: %w", err)
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
