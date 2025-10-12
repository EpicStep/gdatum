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

// ListMultiplayerSummaries ...
func (h *Handlers) ListMultiplayerSummaries(ctx context.Context, params api.ListMultiplayerSummariesParams) ([]api.MultiplayerSummary, error) {
	summary, err := h.repo.ListMultiplayerSummaries(ctx, params.PlayersOrderAsc.Value)
	if err != nil {
		return nil, fmt.Errorf("h.repo.ListMultiplayerSummaries: %w", err)
	}

	return lo.Map(summary, func(multiplayer domain.MultiplayerSummary, _ int) api.MultiplayerSummary {
		return api.MultiplayerSummary{
			Name:         string(multiplayer.Name),
			PlayersCount: multiplayer.PlayersCount,
		}
	}), nil
}

// ListServerSummaries ...
func (h *Handlers) ListServerSummaries(ctx context.Context, params api.ListServerSummariesParams) (api.ListServerSummariesRes, error) {
	servers, err := h.repo.ListServerSummaries(ctx, domain.ListServerSummariesParams{
		Multiplayer:     domain.Multiplayer(params.MultiplayerName),
		IncludeOffline:  params.IncludeOffline.Value,
		PlayersOrderAsc: params.PlayersOrderAsc.Value,
		Limit:           params.Limit.Value,
		Offset:          params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("h.repo.ListServerSummaries: %w", err)
	}

	resp := api.ListServerSummariesOKApplicationJSON(lo.Map(servers, func(server domain.ServerSummary, _ int) api.ServerSummary {
		return api.ServerSummary{
			Host:         server.Host,
			Name:         server.Name,
			PlayersCount: server.PlayersCount,
		}
	}))

	return &resp, nil
}

// GetServer ...
func (h *Handlers) GetServer(ctx context.Context, params api.GetServerParams) (api.GetServerRes, error) {
	server, err := h.repo.GetServer(ctx, domain.Multiplayer(params.MultiplayerName), params.ServerHost)
	if err != nil {
		if errors.Is(err, domain.ErrServerNotFound) {
			return &api.GetServerNotFound{}, nil
		}

		return nil, fmt.Errorf("h.repo.GetServer: %w", err)
	}

	return bindDetailedServer(server), nil
}

// ListServerStatistics ...
func (h *Handlers) ListServerStatistics(ctx context.Context, params api.ListServerStatisticsParams) (api.ListServerStatisticsRes, error) {
	statistics, err := h.repo.ListServerStatistics(ctx, domain.ListServerStatisticsParams{
		Multiplayer: domain.Multiplayer(params.MultiplayerName),
		Host:        params.ServerHost,
		TimeRange: domain.TimeRange{
			From: params.From,
			To:   params.To,
		},
		Precision: precisionToDomain(params.Precision.Value),
	})
	if err != nil {
		if errors.Is(err, domain.ErrServerNotFound) {
			return &api.ListServerStatisticsNotFound{}, nil
		}

		return nil, fmt.Errorf("h.repo.ListServerStatistics: %w", err)
	}

	resp := api.ListServerStatisticsOKApplicationJSON(lo.Map(statistics, func(point domain.ServerStatisticPoint, _ int) api.ServerStatisticPoint {
		return api.ServerStatisticPoint{
			CollectedAt:  point.CollectedAt,
			PlayersCount: point.PlayersCount,
		}
	}))

	return &resp, nil
}

func precisionToDomain(precision api.ListServerStatisticsPrecision) domain.ServerStatisticsPrecision {
	if precision == api.ListServerStatisticsPrecisionPerDay {
		return domain.ServerStatisticsPrecisionPerDay
	}

	return domain.ServerStatisticsPrecisionPerHour
}

func bindDetailedServer(server domain.Server) *api.DetailedServer {
	result := &api.DetailedServer{
		Name: server.Name,
	}

	if server.URL != "" {
		result.URL = api.NewOptString(server.URL)
	}

	if server.Gamemode != "" {
		result.Gamemode = api.NewOptString(server.Gamemode)
	}

	if server.Language != "" {
		result.Language = api.NewOptString(server.Language)
	}

	if server.PlayersCount > 0 {
		result.PlayersCount = api.NewOptInt64(int64(server.PlayersCount))
	}

	return result
}
