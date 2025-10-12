-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_metrics_raw
(
    multiplayer   LowCardinality(String),
    host          String,
    name          String,
    url           String,
    gamemode      String,
    language      String,
    players_count Int32,
    collected_at  Datetime
) ENGINE = Null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_metrics_raw;
-- +goose StatementEnd
