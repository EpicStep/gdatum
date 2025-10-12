-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW servers_online_mv TO servers_online AS
SELECT multiplayer,
       host,
       players_count,
       collected_at
FROM servers_metrics_raw
GROUP BY multiplayer, host, players_count, collected_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_online_mv;
-- +goose StatementEnd
