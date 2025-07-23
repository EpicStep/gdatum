-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_online
(
    multiplayer LowCardinality(String),
    identifier  String,
    players     Int32,
    timestamp   Datetime
) ENGINE = MergeTree()
    ORDER BY (identifier, multiplayer, timestamp)
    PARTITION BY toYYYYMM(timestamp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_online;
-- +goose StatementEnd
