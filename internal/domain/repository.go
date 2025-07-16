// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"context"
	"time"
)

// Repository ...
type Repository interface {
	InsertServers(ctx context.Context, servers []Server) error
	MultiplayersSummary(ctx context.Context, playersOrder Order) ([]MultiplayerSummary, error)
	ServersByMultiplayer(ctx context.Context, filter ServersByMultiplayerFilter) ([]ServerSummary, error)
	ServerByIdentifier(ctx context.Context, multiplayer Multiplayer, identifier string) (Server, error)
	ServerStatistics(ctx context.Context, filter ServerStatisticsFilter) ([]ServerStatistic, error)
}

// Order ...
type Order uint8

const (
	// OrderAsc ...
	OrderAsc Order = iota
	// OrderDesc ...
	OrderDesc
)

// ServersByMultiplayerFilter ...
type ServersByMultiplayerFilter struct {
	Multiplayer    Multiplayer
	Count          int32
	PlayersOrder   Order
	IncludeOffline bool
}

// ServerStatisticsFilterPrecision ...
type ServerStatisticsFilterPrecision uint8

const (
	// ServerStatisticsFilterPrecisionPerHour ...
	ServerStatisticsFilterPrecisionPerHour ServerStatisticsFilterPrecision = iota
	// ServerStatisticsFilterPrecisionPerDay ...
	ServerStatisticsFilterPrecisionPerDay
)

// ServerStatisticsFilter ...
type ServerStatisticsFilter struct {
	Multiplayer Multiplayer
	Identifier  string
	TimeOrder   Order
	Count       int32
	LastSeen    time.Time
	Precision   ServerStatisticsFilterPrecision
}
