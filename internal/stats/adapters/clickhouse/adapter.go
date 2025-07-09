package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/EpicStep/gdatum/internal/stats/domain"
	"github.com/EpicStep/gdatum/internal/stats/repository"
	"github.com/EpicStep/gdatum/internal/stats/repository/clickhouse"
	"github.com/samber/lo"
)

type Adapter struct {
	store *clickhouse.Store
}

func New(store *clickhouse.Store) *Adapter {
	return &Adapter{
		store: store,
	}
}

func (a *Adapter) InsertServers(ctx context.Context, servers []domain.Server) error {
	chServers := lo.Map(servers, func(srv domain.Server, index int) clickhouse.Server {
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

func (a *Adapter) GetMultiplayersSummary(ctx context.Context, playersOrder repository.Order) ([]domain.MultiplayerSummary, error) {
	summaries, err := a.store.GetMultiplayersSummary(ctx, playersOrder)
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

func (a *Adapter) GetServerByID(ctx context.Context, multiplayer string, identifier string) (domain.Server, error) {
	chServer, err := a.store.GetServerByID(ctx, multiplayer, identifier)
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

func (a *Adapter) GetServersByMultiplayer(ctx context.Context, filter repository.GetServersByMultiplayerFilter) ([]domain.ServerSummary, error) {
	servers, err := a.store.GetServersByMultiplayer(ctx, filter)
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

func (a *Adapter) GetServerStatistics(ctx context.Context, filter repository.GetServerStatisticsFilter) ([]domain.ServerStatistic, error) {
	statistics, err := a.store.GetServerStatistics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("a.store.GetServerStatistics: %w", err)
	}

	if len(statistics) == 0 {
		return nil, domain.ErrServerNotFound
	}

	return lo.Map(statistics, func(stat clickhouse.ServerStatistic, _ int) domain.ServerStatistic {
		return domain.ServerStatistic{
			Players:   stat.Players,
			Timestamp: stat.Timestamp,
		}
	}), nil
}
