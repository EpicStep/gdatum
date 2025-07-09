-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers_metrics_raw
(
    multiplayer LowCardinality(String),
    identifier  String,
    name        String,
    url         String,
    gamemode    String,
    lang        String,
    players     Int32,
    timestamp   Datetime
) ENGINE = Null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_metrics_raw;
-- +goose StatementEnd
