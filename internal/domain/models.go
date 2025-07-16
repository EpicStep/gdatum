// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"time"
)

// Multiplayer ...
type Multiplayer string

const (
	// MultiplayerRagemp ...
	MultiplayerRagemp = "ragemp"
)

// Server ...
type Server struct {
	Multiplayer Multiplayer
	Identifier  string
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
	Identifier string
	Name       string
	Players    int32
}
