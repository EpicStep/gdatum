// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package api

import (
	"context"
	"fmt"
	"github.com/EpicStep/gdatum/internal/stats/domain"
	"github.com/EpicStep/gdatum/internal/stats/repository"
	"github.com/EpicStep/gdatum/pkg/api"
	"github.com/go-faster/errors"
	"github.com/samber/lo"
)

var _ api.Handler = (*Handlers)(nil)

// Handlers ...
type Handlers struct {
	repo repository.Repository
}

// New returns new Handlers.
func New(repo repository.Repository) *Handlers {
	return &Handlers{
		repo: repo,
	}
}

func (h *Handlers) GetMultiplayersSummary(ctx context.Context, params api.GetMultiplayersSummaryParams) ([]api.MultiplayerSummary, error) {
	summary, err := h.repo.GetMultiplayersSummary(ctx, orderToRepository(string(params.PlayersOrder.Value)))
	if err != nil {
		return nil, fmt.Errorf("h.repo.GetMultiplayersSummary: %w", err)
	}

	return lo.Map(summary, func(multiplayer domain.MultiplayerSummary, _ int) api.MultiplayerSummary {
		return api.MultiplayerSummary{
			Name:    string(multiplayer.Name),
			Players: multiplayer.Players,
		}
	}), nil
}

func (h *Handlers) GetServersByMultiplayer(ctx context.Context, params api.GetServersByMultiplayerParams) (api.GetServersByMultiplayerRes, error) {
	servers, err := h.repo.GetServersByMultiplayer(ctx, repository.GetServersByMultiplayerFilter{
		Multiplayer:    domain.Multiplayer(params.MultiplayerName),
		Count:          params.Count.Value,
		PlayersOrder:   orderToRepository(string(params.PlayersOrder.Value)),
		IncludeOffline: params.IncludeOffline.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("h.repo.GetServersByMultiplayer: %w", err)
	}

	resp := api.GetServersByMultiplayerOKApplicationJSON(lo.Map(servers, func(server domain.ServerSummary, _ int) api.ServerSummary {
		return api.ServerSummary{
			ID:      server.Identifier,
			Name:    server.Name,
			Players: server.Players,
		}
	}))

	return &resp, nil
}

func (h *Handlers) GetServerByID(ctx context.Context, params api.GetServerByIDParams) (api.GetServerByIDRes, error) {
	server, err := h.repo.GetServerByID(ctx, params.MultiplayerName, params.ServerID)
	if err != nil {
		if errors.Is(err, domain.ErrServerNotFound) {
			return &api.GetServerByIDNotFound{}, nil
		}

		return nil, fmt.Errorf("h.repo.GetServerByID: %w", err)
	}

	return &api.GetServerByIDOK{
		Name:        server.Name,
		URL:         api.NewOptString(server.URL),
		Gamemode:    api.NewOptString(server.Gamemode),
		Lang:        api.NewOptString(server.Lang),
		Players:     api.NewOptInt64(int64(server.Players)),
		CollectedAt: api.NewOptDateTime(server.CollectedAt),
	}, nil
}

func (h *Handlers) GetServerStatisticsByID(ctx context.Context, params api.GetServerStatisticsByIDParams) (api.GetServerStatisticsByIDRes, error) {
	statistics, err := h.repo.GetServerStatistics(ctx, repository.GetServerStatisticsFilter{
		Multiplayer: domain.Multiplayer(params.MultiplayerName),
		Identifier:  params.ServerID,
		TimeOrder:   orderToRepository(string(params.TimeOrder.Value)),
		Count:       params.Count.Value,
		LastSeen:    params.LastSeen.Value,
		Precision:   precisionToRepository(params.Precision.Value),
	})
	if err != nil {
		if errors.Is(err, domain.ErrServerNotFound) {
			return &api.GetServerStatisticsByIDNotFound{}, nil
		}

		return nil, fmt.Errorf("h.repo.GetServerStatisticsByID: %w", err)
	}

	resp := api.GetServerStatisticsByIDOKApplicationJSON(lo.Map(statistics, func(stat domain.ServerStatistic, _ int) api.ServerStatistic {
		return api.ServerStatistic{
			Timestamp: stat.Timestamp,
			Players:   stat.Players,
		}
	}))

	return &resp, nil
}

func orderToRepository(order string) repository.Order {
	switch order {
	case "asc":
		return repository.OrderAsc
	}

	return repository.OrderDesc
}

func precisionToRepository(precision api.GetServerStatisticsByIDPrecision) repository.GetServerStatisticsFilterPrecision {
	switch precision {
	case api.GetServerStatisticsByIDPrecisionPerDay:
		return repository.GetServerStatisticsFilterPrecisionPerDay
	}

	return repository.GetServerStatisticsFilterPrecisionPerHour
}
