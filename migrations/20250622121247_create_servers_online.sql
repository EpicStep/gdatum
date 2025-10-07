-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_online
(
    multiplayer  LowCardinality(String),
    id            String,
    players_count Int32 CODEC(T64, ZSTD),
    collected_at  Datetime CODEC(DoubleDelta, ZSTD)
) ENGINE = MergeTree()
      ORDER BY (id, multiplayer, collected_at)
      PARTITION BY toYYYYMM(collected_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_online;
-- +goose StatementEnd
