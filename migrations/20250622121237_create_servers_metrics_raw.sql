-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_metrics_raw
(
    multiplayer   LowCardinality(String),
    id            String, -- it’s almost an IP address, but sometimes it may not be
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
