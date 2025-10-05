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

// Server ...
type Server struct {
	Multiplayer Multiplayer
	ID          string
	Name        string
	URL         string
	Gamemode    string
	Lang        string
	Players     int32
	CollectedAt time.Time
}

// MultiplayerSummary ...
type MultiplayerSummary struct {
	Name    Multiplayer
	Players int64
}

// ServerStatistic ..
type ServerStatistic struct {
	Players int32
	At      time.Time
}

// ServerSummary ...
type ServerSummary struct {
	ID      string
	Name    string
	Players int32
}
