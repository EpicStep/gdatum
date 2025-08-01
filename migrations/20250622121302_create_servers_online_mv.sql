-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW servers_online_mv TO servers_online AS
SELECT multiplayer,
       identifier,
       players,
       timestamp
FROM servers_metrics_raw
GROUP BY multiplayer, identifier, players, timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_online_mv;
-- +goose StatementEnd
