CREATE MATERIALIZED VIEW nats.logs_view
            ENGINE = Memory
AS
SELECT * FROM nats.logs
SETTINGS
    stream_like_engine_allow_direct_select = 1;