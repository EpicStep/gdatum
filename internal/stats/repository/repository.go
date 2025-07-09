package repository

import (
	"context"
	"github.com/EpicStep/gdatum/internal/stats/domain"
	"time"
)

type Repository interface {
	InsertServers(ctx context.Context, servers []domain.Server) error
	GetMultiplayersSummary(ctx context.Context, playersOrder Order) ([]domain.MultiplayerSummary, error)
	GetServersByMultiplayer(ctx context.Context, filter GetServersByMultiplayerFilter) ([]domain.ServerSummary, error)
	GetServerByID(ctx context.Context, multiplayer string, identifier string) (domain.Server, error)
	GetServerStatistics(ctx context.Context, filter GetServerStatisticsFilter) ([]domain.ServerStatistic, error)
}

type Order uint8

const (
	OrderAsc Order = iota
	OrderDesc
)

type GetServersByMultiplayerFilter struct {
	Multiplayer    domain.Multiplayer
	Count          int32
	PlayersOrder   Order
	IncludeOffline bool
}

type GetServerStatisticsFilterPrecision uint8

const (
	GetServerStatisticsFilterPrecisionPerHour GetServerStatisticsFilterPrecision = iota
	GetServerStatisticsFilterPrecisionPerDay
)

type GetServerStatisticsFilter struct {
	Multiplayer domain.Multiplayer
	Identifier  string
	TimeOrder   Order
	Count       int32
	LastSeen    time.Time
	Precision   GetServerStatisticsFilterPrecision
}
