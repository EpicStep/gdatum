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
		Cols(multiplayerColumnName, hostColumnName, nameColumnName, languageColumnName, gamemodeColumnName, urlColumnName, playersCountColumnName, collectedAtColumnName)

	sqlRaw, _ := sql.Build(ib)

	batch, err := s.db.PrepareBatch(ctx, sqlRaw)
	if err != nil {
		return fmt.Errorf("s.db.PrepareBatch: %w", err)
	}

	for _, server := range servers {
		err = batch.Append(
			server.Multiplayer,
			server.Host,
			server.Name,
			server.Language,
			server.Gamemode,
			server.URL,
			server.PlayersCount,
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

// ListMultiplayerSummaries ...
func (s *Store) ListMultiplayerSummaries(ctx context.Context, playersOrderAsc bool) ([]MultiplayerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversOnlineTableName).
		Select(multiplayerColumnName, sb.As(wrapColumn("sum", playersCountColumnName), playersCountColumnName)).
		Where(fmt.Sprintf("%s = toStartOfHour(now())", collectedAtColumnName)).
		GroupBy(multiplayerColumnName)

	if !playersOrderAsc {
		sb = sb.OrderBy(playersCountColumnName + " DESC")
	}

	sqlRaw, args := sb.Build()

	var result []MultiplayerSummary
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("s.db.Select: %w", err)
	}

	return result, nil
}

// ListServerSummaries ...
func (s *Store) ListServerSummaries(ctx context.Context, params domain.ListServerSummariesParams) ([]ServerSummary, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(hostColumnName, nameColumnName, playersCountColumnName).
		Where(sb.Equal(multiplayerColumnName, string(params.Multiplayer))).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, hostColumnName, serversOnlineTableName, hostColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, collectedAtColumnName),
		).
		OrderBy(collectedAtColumnName + " DESC").
		Limit(int(params.Limit)).
		Offset(int(params.Offset))

	if !params.PlayersOrderAsc {
		sb = sb.OrderBy(playersCountColumnName + " DESC")
	}

	sqlRaw, args := sb.Build()

	var result []ServerSummary
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("s.db.Select: %w", err)
	}

	return result, nil
}

// GetServer ...
func (s *Store) GetServer(ctx context.Context, multiplayer domain.Multiplayer, host string) (Server, error) {
	sb := sqlbuilder.NewSelectBuilder()

	sb = sb.From(serversInfoTableName).
		Select(multiplayerColumnName, hostColumnName, nameColumnName, languageColumnName, gamemodeColumnName, urlColumnName, playersCountColumnName, collectedAtColumnName).
		Where(
			sb.And(
				sb.Equal(multiplayerColumnName, multiplayer),
				sb.Equal(hostColumnName, host),
			),
		).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			serversOnlineTableName,
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, multiplayerColumnName, serversOnlineTableName, multiplayerColumnName),
			fmt.Sprintf("%s.%s = %s.%s", serversInfoTableName, hostColumnName, serversOnlineTableName, hostColumnName),
			fmt.Sprintf("%s.%s = toStartOfHour(now())", serversOnlineTableName, collectedAtColumnName),
		).OrderBy(collectedAtColumnName + " DESC").Limit(1)

	sqlRaw, args := sb.Build()

	var srv Server
	if err := s.db.QueryRow(ctx, sqlRaw, args...).ScanStruct(&srv); err != nil {
		return Server{}, fmt.Errorf("s.db.QueryRow: %w", err)
	}

	return srv, nil
}

// ListServerStatistics ...
func (s *Store) ListServerStatistics(ctx context.Context, params domain.ListServerStatisticsParams) ([]ServerStatisticPoint, error) {
	sb := sqlbuilder.NewSelectBuilder()

	timeSelect := wrapColumn("toStartOfHour", collectedAtColumnName)
	if params.Precision == domain.ServerStatisticsPrecisionPerDay {
		timeSelect = wrapColumn("toStartOfDay", collectedAtColumnName)
	}

	sb = sb.
		From(serversOnlineTableName).
		Select(
			sb.As(timeSelect, collectedAtColumnName),
			sb.As(wrapColumn("toInt32", wrapColumn("avg", playersCountColumnName)), playersCountColumnName),
		).
		Where(
			sb.Equal(multiplayerColumnName, string(params.Multiplayer)),
			sb.Equal(hostColumnName, params.Host),
			sb.GreaterThan(collectedAtColumnName, params.TimeRange.From),
			sb.LessThan(collectedAtColumnName, params.TimeRange.To),
		).
		GroupBy(collectedAtColumnName).
		OrderBy(collectedAtColumnName + " DESC")

	sqlRaw, args := sql.Build(sb)

	var result []ServerStatisticPoint
	if err := s.db.Select(ctx, &result, sqlRaw, args...); err != nil {
		return nil, fmt.Errorf("s.db.Select: %w", err)
	}

	return result, nil
}

func wrapColumn(wrapper, columnName string) string {
	return wrapper + "(" + columnName + ")"
}
