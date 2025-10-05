// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package api

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/samber/lo"

	"github.com/EpicStep/gdatum/internal/domain"
	"github.com/EpicStep/gdatum/pkg/api"
)

var _ api.Handler = (*Handlers)(nil)

// Handlers ...
type Handlers struct {
	repo domain.Repository
}

// New returns new Handlers.
func New(repo domain.Repository) *Handlers {
	return &Handlers{
		repo: repo,
	}
}

// GetMultiplayersSummary ...
func (h *Handlers) GetMultiplayersSummary(ctx context.Context, params api.GetMultiplayersSummaryParams) ([]api.MultiplayerSummary, error) {
	summary, err := h.repo.GetMultiplayersSummary(ctx, orderToDomain(params.PlayersOrder.Value))
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

// GetServersByMultiplayer ...
func (h *Handlers) GetServersByMultiplayer(ctx context.Context, params api.GetServersByMultiplayerParams) (api.GetServersByMultiplayerRes, error) {
	servers, err := h.repo.GetServersByMultiplayer(ctx, domain.ServersByMultiplayerFilter{
		Multiplayer:    domain.Multiplayer(params.MultiplayerName),
		Count:          params.Count.Value,
		PlayersOrder:   orderToDomain(params.PlayersOrder.Value),
		IncludeOffline: params.IncludeOffline.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("h.repo.GetServersByMultiplayer: %w", err)
	}

	resp := api.GetServersByMultiplayerOKApplicationJSON(lo.Map(servers, func(server domain.ServerSummary, _ int) api.ServerSummary {
		return api.ServerSummary{
			ID:      server.ID,
			Name:    server.Name,
			Players: server.Players,
		}
	}))

	return &resp, nil
}

// GetServerByID ...
func (h *Handlers) GetServerByID(ctx context.Context, params api.GetServerByIDParams) (api.GetServerByIDRes, error) {
	server, err := h.repo.GetServerByID(ctx, domain.Multiplayer(params.MultiplayerName), params.ServerID)
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

// GetServerStatisticsByID ...
func (h *Handlers) GetServerStatisticsByID(ctx context.Context, params api.GetServerStatisticsByIDParams) (api.GetServerStatisticsByIDRes, error) {
	statistics, err := h.repo.GetServerStatistics(ctx, domain.ServerStatisticsFilter{
		Multiplayer: domain.Multiplayer(params.MultiplayerName),
		ID:          params.ServerID,
		TimeOrder:   orderToDomain(params.TimeOrder.Value),
		Count:       params.Count.Value,
		LastSeen:    params.LastSeen.Value,
		Precision:   precisionToDomain(params.Precision.Value),
	})
	if err != nil {
		if errors.Is(err, domain.ErrServerNotFound) {
			return &api.GetServerStatisticsByIDNotFound{}, nil
		}

		return nil, fmt.Errorf("h.repo.GetServerStatisticsByID: %w", err)
	}

	resp := api.GetServerStatisticsByIDOKApplicationJSON(lo.Map(statistics, func(stat domain.ServerStatistic, _ int) api.ServerStatistic {
		return api.ServerStatistic{
			At:      stat.At,
			Players: stat.Players,
		}
	}))

	return &resp, nil
}

func orderToDomain(order api.Order) domain.Order {
	if order == api.OrderAsc {
		return domain.OrderAsc
	}

	return domain.OrderDesc
}

func precisionToDomain(precision api.GetServerStatisticsByIDPrecision) domain.ServerStatisticsFilterPrecision {
	if precision == api.GetServerStatisticsByIDPrecisionPerDay {
		return domain.ServerStatisticsFilterPrecisionPerDay
	}

	return domain.ServerStatisticsFilterPrecisionPerHour
}
