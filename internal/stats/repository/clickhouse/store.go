// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/EpicStep/gdatum/internal/stats/repository"
	"github.com/EpicStep/gdatum/internal/utils/sql"
	"github.com/huandu/go-sqlbuilder"
)

// Store ...
type Store struct {
	db driver.Conn
}

// New ...
func New(db driver.Conn) *Store {
	return &Store{
		db: db,
	}
}

// InsertServers ...
func (s *Store) InsertServers(ctx context.Context, servers []Server) error {
	if len(servers) == 0 {
		return nil
	}

	ib := sqlbuilder.
		NewInsertBuilder().
		InsertInto(serversMetricsRawTableName).
		Cols(multiplayerColumnName, identifierColumnName, nameColumnName, langColumnName, gamemodeColumnName, urlColumnName, playersColumnName, timestampColumnName)

	sqlRaw, _ := sql.Build(ib)

	batch, err := s.db.PrepareBatch(ctx, sqlRaw)
	if err != nil {
		return fmt.Errorf("r.db.PrepareBatch: %w", err)
	}

	for _, server := range servers {
		err = batch.Append(
			server.Multiplayer,
			server.Identifier,
			server.Name,
			server.Lang,
			server.Gamemode,
			server.URL,
			server.Players,
			server.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("batch.Append: %w", err)
		}
	}

	if err = batch.Send(); err != nil {
		return fmt.Errorf("batch.Send: %w", err)
	}

	return nil
}

func (s *Store) GetMultiplayersSummary(ctx context.Context, playersOrder repository.Order) ([]MultiplayerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversOnlineTableName).
		Select(multiplayerColumnName, sb.As(wrapColumn("sum", playersColumnName), onlineColumnName)).
		Where(fmt.Sprintf("%s = toStartOfHour(now())", timestampColumnName)).
		GroupBy(multiplayerColumnName)

	if playersOrder == repository.OrderAsc {
		sb = sb.OrderBy(onlineColumnName + " ASC")
	} else {
		sb = sb.OrderBy(onlineColumnName + " DESC")
	}

	sqlRaw, args := sb.Build()

	var result []MultiplayerSummary
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("r.db.Select: %w", err)
	}

	return result, nil
}

func (s *Store) GetServersByMultiplayer(ctx context.Context, filter repository.GetServersByMultiplayerFilter) ([]ServerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(identifierColumnName, nameColumnName, playersColumnName).
		Where(sb.Equal(multiplayerColumnName, string(filter.Multiplayer))).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, identifierColumnName, serversOnlineTableName, identifierColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, timestampColumnName),
		).OrderBy(timestampColumnName + " DESC")

	if filter.PlayersOrder == repository.OrderAsc {
		sb = sb.OrderBy(playersColumnName + " ASC")
	} else {
		sb = sb.OrderBy(playersColumnName + " DESC")
	}

	if filter.Count >= 0 {
		sb = sb.Limit(int(filter.Count))
	}

	sqlRaw, args := sb.Build()

	var result []ServerSummary
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("r.db.Select: %w", err)
	}

	return result, nil
}

func (s *Store) GetServerByID(ctx context.Context, multiplayer string, identifier string) (Server, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(multiplayerColumnName, identifierColumnName, nameColumnName, langColumnName, gamemodeColumnName, urlColumnName, playersColumnName, timestampColumnName).
		Where(
			sb.And(
				sb.Equal(multiplayerColumnName, multiplayer),
				sb.Equal(identifierColumnName, identifier),
			),
		).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, identifierColumnName, serversOnlineTableName, identifierColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, timestampColumnName),
		).OrderBy(timestampColumnName + " DESC").Limit(1)

	sqlRaw, args := sb.Build()

	var srv Server
	if err := s.db.QueryRow(ctx, sqlRaw, args...).ScanStruct(&srv); err != nil {
		return Server{}, fmt.Errorf("r.db.QueryRow: %w", err)
	}

	return srv, nil
}

func (s *Store) GetServerStatistics(ctx context.Context, filter repository.GetServerStatisticsFilter) ([]ServerStatistic, error) {
	sb := sqlbuilder.NewSelectBuilder()

	timeSelect := wrapColumn("toStartOfHour", timestampColumnName)
	if filter.Precision == repository.GetServerStatisticsFilterPrecisionPerDay {
		timeSelect = wrapColumn("toStartOfDay", timestampColumnName)
	}

	sb = sb.
		From(serversOnlineTableName).
		Select(
			sb.As(timeSelect, timestampColumnName),
			sb.As(wrapColumn("toInt32", wrapColumn("avg", playersColumnName)), playersColumnName),
		).
		Where(
			sb.Equal(multiplayerColumnName, string(filter.Multiplayer)),
			sb.Equal(identifierColumnName, filter.Identifier),
		).
		GroupBy(timestampColumnName).
		Limit(int(filter.Count))

	if !filter.LastSeen.IsZero() {
		if filter.TimeOrder == repository.OrderAsc {
			sb = sb.Where(sb.GreaterThan(timestampColumnName, filter.LastSeen))
		} else {
			sb = sb.Where(sb.LessThan(timestampColumnName, filter.LastSeen))
		}
	}

	if filter.TimeOrder == repository.OrderAsc {
		sb = sb.OrderBy(timestampColumnName + " ASC")
	} else {
		sb = sb.OrderBy(timestampColumnName + " DESC")
	}

	sqlRaw, args := sql.Build(sb)

	var result []ServerStatistic
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("r.db.Select: %w", err)
	}

	return result, nil
}

func wrapColumn(wrapper, columnName string) string {
	return wrapper + "(" + columnName + ")"
}
