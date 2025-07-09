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
)

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

type MultiplayerSummary struct {
	Multiplayer string `ch:"multiplayer"`
	Players     int64  `ch:"online"`
}

type ServerStatistic struct {
	Players   int32     `ch:"players"`
	Timestamp time.Time `ch:"timestamp"`
}

type ServerSummary struct {
	Identifier string `ch:"identifier"`
	Name       string `ch:"name"`
	Players    int32  `ch:"players"`
}
