// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/huandu/go-sqlbuilder"

	"github.com/EpicStep/gdatum/internal/domain"
	"github.com/EpicStep/gdatum/internal/utils/sql"
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
		Cols(multiplayerColumnName, idColumnName, nameColumnName, langColumnName, gamemodeColumnName, urlColumnName, playersColumnName, collectedAtColumnName)

	sqlRaw, _ := sql.Build(ib)

	batch, err := s.db.PrepareBatch(ctx, sqlRaw)
	if err != nil {
		return fmt.Errorf("s.db.PrepareBatch: %w", err)
	}

	for _, server := range servers {
		err = batch.Append(
			server.Multiplayer,
			server.ID,
			server.Name,
			server.Lang,
			server.Gamemode,
			server.URL,
			server.Players,
			server.CollectedAt,
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

// GetMultiplayersSummary ...
func (s *Store) GetMultiplayersSummary(ctx context.Context, playersOrder domain.Order) ([]MultiplayerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversOnlineTableName).
		Select(multiplayerColumnName, sb.As(wrapColumn("sum", playersColumnName), onlineColumnName)).
		Where(fmt.Sprintf("%s = toStartOfHour(now())", collectedAtColumnName)).
		GroupBy(multiplayerColumnName)

	if playersOrder == domain.OrderAsc {
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

// GetServersByMultiplayer ...
func (s *Store) GetServersByMultiplayer(ctx context.Context, filter domain.ServersByMultiplayerFilter) ([]ServerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(idColumnName, nameColumnName, playersColumnName).
		Where(sb.Equal(multiplayerColumnName, string(filter.Multiplayer))).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, idColumnName, serversOnlineTableName, idColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, collectedAtColumnName),
		).OrderBy(collectedAtColumnName + " DESC")

	if filter.PlayersOrder == domain.OrderAsc {
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

// GetServerByID ...
func (s *Store) GetServerByID(ctx context.Context, multiplayer domain.Multiplayer, id string) (Server, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(multiplayerColumnName, idColumnName, nameColumnName, langColumnName, gamemodeColumnName, urlColumnName, playersColumnName, collectedAtColumnName).
		Where(
			sb.And(
				sb.Equal(multiplayerColumnName, multiplayer),
				sb.Equal(idColumnName, id),
			),
		).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, idColumnName, serversOnlineTableName, idColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, collectedAtColumnName),
		).OrderBy(collectedAtColumnName + " DESC").Limit(1)

	sqlRaw, args := sb.Build()

	var srv Server
	if err := s.db.QueryRow(ctx, sqlRaw, args...).ScanStruct(&srv); err != nil {
		return Server{}, fmt.Errorf("r.db.QueryRow: %w", err)
	}

	return srv, nil
}

// GetServerStatistics ...
func (s *Store) GetServerStatistics(ctx context.Context, filter domain.ServerStatisticsFilter) ([]ServerStatistic, error) {
	sb := sqlbuilder.NewSelectBuilder()

	timeSelect := wrapColumn("toStartOfHour", collectedAtColumnName)
	if filter.Precision == domain.ServerStatisticsFilterPrecisionPerDay {
		timeSelect = wrapColumn("toStartOfDay", collectedAtColumnName)
	}

	sb = sb.
		From(serversOnlineTableName).
		Select(
			sb.As(timeSelect, atColumnName),
			sb.As(wrapColumn("toInt32", wrapColumn("avg", playersColumnName)), playersColumnName),
		).
		Where(
			sb.Equal(multiplayerColumnName, string(filter.Multiplayer)),
			sb.Equal(idColumnName, filter.ID),
		).
		GroupBy(collectedAtColumnName).
		Limit(int(filter.Count))

	if !filter.LastSeen.IsZero() {
		if filter.TimeOrder == domain.OrderAsc {
			sb = sb.Where(sb.GreaterThan(collectedAtColumnName, filter.LastSeen))
		} else {
			sb = sb.Where(sb.LessThan(collectedAtColumnName, filter.LastSeen))
		}
	}

	if filter.TimeOrder == domain.OrderAsc {
		sb = sb.OrderBy(collectedAtColumnName + " ASC")
	} else {
		sb = sb.OrderBy(collectedAtColumnName + " DESC")
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
