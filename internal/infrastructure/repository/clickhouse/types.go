// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import "time"

const (
	serversMetricsRawTableName = "servers_metrics_raw"
	serversInfoTableName       = "servers_info"
	serversOnlineTableName     = "servers_online"

	multiplayerColumnName = "multiplayer"
	idColumnName          = "id"
	nameColumnName        = "name"
	langColumnName        = "lang"
	gamemodeColumnName    = "gamemode"
	urlColumnName         = "url"
	playersColumnName     = "players"
	collectedAtColumnName = "collected_at"
	onlineColumnName      = "online"
	atColumnName          = "at"
)

// Server ...
type Server struct {
	Multiplayer string    `ch:"multiplayer"`
	ID          string    `ch:"id"`
	Name        string    `ch:"name"`
	URL         string    `ch:"url"`
	Gamemode    string    `ch:"gamemode"`
	Lang        string    `ch:"lang"`
	Players     int32     `ch:"players"`
	CollectedAt time.Time `ch:"collected_at"`
}

// MultiplayerSummary ...
type MultiplayerSummary struct {
	Multiplayer string `ch:"multiplayer"`
	Players     int64  `ch:"online"`
}

// ServerStatistic ...
type ServerStatistic struct {
	Players int32     `ch:"players"`
	At      time.Time `ch:"at"`
}

// ServerSummary ...
type ServerSummary struct {
	ID      string `ch:"id"`
	Name    string `ch:"name"`
	Players int32  `ch:"players"`
}
