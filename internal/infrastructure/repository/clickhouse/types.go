// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import "time"

const (
	serversMetricsRawTableName = "servers_metrics_raw"
	serversInfoTableName       = "servers_info"
	serversOnlineTableName     = "servers_online"

	multiplayerColumnName = "multiplayer"
	identifierColumnName  = "identifier"
	nameColumnName        = "name"
	langColumnName        = "lang"
	gamemodeColumnName    = "gamemode"
	urlColumnName         = "url"
	playersColumnName     = "players"
	timestampColumnName   = "timestamp"
	onlineColumnName      = "online"
	atColumnName          = "at"
)

// Server ...
type Server struct {
	Multiplayer string    `ch:"multiplayer"`
	Identifier  string    `ch:"identifier"`
	Name        string    `ch:"name"`
	URL         string    `ch:"url"`
	Gamemode    string    `ch:"gamemode"`
	Lang        string    `ch:"lang"`
	Players     int32     `ch:"players"`
	Timestamp   time.Time `ch:"timestamp"`
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
	Identifier string `ch:"identifier"`
	Name       string `ch:"name"`
	Players    int32  `ch:"players"`
}
