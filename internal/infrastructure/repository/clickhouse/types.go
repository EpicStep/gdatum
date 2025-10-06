// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import "time"

const (
	serversMetricsRawTableName = "servers_metrics_raw"
	serversInfoTableName       = "servers_info"
	serversOnlineTableName     = "servers_online"

	multiplayerColumnName  = "multiplayer"
	idColumnName           = "id"
	nameColumnName         = "name"
	languageColumnName     = "language"
	gamemodeColumnName     = "gamemode"
	urlColumnName          = "url"
	playersCountColumnName = "players_count"
	collectedAtColumnName  = "collected_at"
)

// Server ...
type Server struct {
	Multiplayer  string    `ch:"multiplayer"`
	ID           string    `ch:"id"`
	Name         string    `ch:"name"`
	URL          string    `ch:"url"`
	Gamemode     string    `ch:"gamemode"`
	Language     string    `ch:"language"`
	PlayersCount int32     `ch:"players_count"`
	CollectedAt  time.Time `ch:"collected_at"`
}

// MultiplayerSummary ...
type MultiplayerSummary struct {
	Multiplayer  string `ch:"multiplayer"`
	PlayersCount int64  `ch:"players_count"`
}

// ServerStatisticPoint ...
type ServerStatisticPoint struct {
	PlayersCount int32     `ch:"players_count"`
	CollectedAt  time.Time `ch:"collected_at"`
}

// ServerSummary ...
type ServerSummary struct {
	ID           string `ch:"id"`
	Name         string `ch:"name"`
	PlayersCount int32  `ch:"players_count"`
}
