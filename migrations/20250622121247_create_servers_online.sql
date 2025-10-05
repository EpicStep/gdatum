-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_online
(
    multiplayer LowCardinality(String),
    id          String,
    players     Int32,
    collected_at   Datetime
) ENGINE = MergeTree()
    ORDER BY (id, multiplayer, collected_at)
    PARTITION BY toYYYYMM(collected_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_online;
-- +goose StatementEnd
