// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"time"
)

// Multiplayer is an alias that represents supported multiplayer's.
type Multiplayer string

const (
	// MultiplayerRagemp ...
	MultiplayerRagemp = "ragemp"
)

// MultiplayerSummary ...
type MultiplayerSummary struct {
	Name         Multiplayer
	PlayersCount int64
}

// ServerStatisticPoint represents a single data point in a graph.
type ServerStatisticPoint struct {
	PlayersCount int32
	CollectedAt  time.Time
}

// Server ...
type Server struct {
	Multiplayer  Multiplayer
	Host         string
	Name         string
	URL          string
	Gamemode     string
	Language     string
	PlayersCount int32
	CollectedAt  time.Time
}

// ServerSummary ...
type ServerSummary struct {
	Host         string
	Name         string
	PlayersCount int32
}
