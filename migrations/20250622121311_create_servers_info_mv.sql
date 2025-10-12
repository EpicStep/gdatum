-- +goose Up
-- +goose StatementBegin
CREATE
MATERIALIZED VIEW servers_info_mv TO servers_info AS
SELECT multiplayer,
       host,
       name,
       url,
       gamemode,
       language,
       collected_at
FROM servers_metrics_raw
GROUP BY multiplayer, host, name, url, gamemode, language, collected_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE servers_info_mv;
-- +goose StatementEnd
