// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Repository ...
type Repository interface {
	InsertServers(ctx context.Context, servers []Server) error
	ListMultiplayerSummaries(ctx context.Context, playersOrderAsc bool) ([]MultiplayerSummary, error)
	ListServerSummaries(ctx context.Context, params ListServerSummariesParams) ([]ServerSummary, error)
	GetServer(ctx context.Context, multiplayer Multiplayer, host string) (Server, error)
	ListServerStatistics(ctx context.Context, params ListServerStatisticsParams) ([]ServerStatisticPoint, error)
}

const (
	serverStatisticsMaxTimeRangeDelta = time.Hour * 24 * 30 // 30 days
)

var (
	errBadLimit  = errors.New("limit must be greater or equal than zero")
	errBadOffset = errors.New("offset must be greater than zero")
)

// ListServerSummariesParams ...
type ListServerSummariesParams struct {
	Multiplayer     Multiplayer
	IncludeOffline  bool
	PlayersOrderAsc bool
	Limit           int32
	Offset          int32
}

func (s ListServerSummariesParams) Validate() error {
	if s.Limit <= 0 {
		return errBadLimit
	}

	if s.Offset < 0 {
		return errBadOffset
	}

	return nil
}

// ServerStatisticsPrecision ...
type ServerStatisticsPrecision uint8

const (
	// ServerStatisticsPrecisionPerHour ...
	ServerStatisticsPrecisionPerHour ServerStatisticsPrecision = iota
	// ServerStatisticsPrecisionPerDay ...
	ServerStatisticsPrecisionPerDay
)

// TimeRange is a type that represents a range of time.
type TimeRange struct {
	From time.Time
	To   time.Time
}

var (
	errIncorrectTimeRange     = errors.New("time range is incorrect, 'From' must be greater than 'To'")
	errTimeRangeDeltaOverflow = errors.New("delta bigger then max")
)

func (t TimeRange) Validate(maxDelta time.Duration) error {
	if (t.From.IsZero() || t.To.IsZero()) || t.To.Before(t.From) {
		return errIncorrectTimeRange
	}

	if maxDelta <= 0 {
		return nil
	}

	if t.To.Sub(t.From) > maxDelta {
		return errTimeRangeDeltaOverflow
	}

	return nil
}

// ListServerStatisticsParams ...
type ListServerStatisticsParams struct {
	Multiplayer Multiplayer
	Host        string
	TimeRange   TimeRange
	Precision   ServerStatisticsPrecision
}

func (s ListServerStatisticsParams) Validate() error {
	if err := s.TimeRange.Validate(serverStatisticsMaxTimeRangeDelta); err != nil {
		return fmt.Errorf("s.TimeRange.Validate: %w", err)
	}

	return nil
}
