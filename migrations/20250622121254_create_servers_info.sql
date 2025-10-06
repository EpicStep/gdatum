-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_info
(
    multiplayer  LowCardinality(String),
    id           String,
    name         String,
    url          String,
    gamemode     String,
    language     String,
    collected_at Datetime
) ENGINE = ReplacingMergeTree()
      ORDER BY (id, multiplayer)
      PARTITION BY toYYYYMM(collected_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_info;
-- +goose StatementEnd
