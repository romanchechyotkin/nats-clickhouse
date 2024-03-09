CREATE TABLE IF NOT EXISTS nats.logs
(
    id UUID,
    `level` Enum8('error' = 0, 'warn' = 1, 'debug' = 2, 'info' = 3),
    `text` String
) ENGINE = NATS()
    SETTINGS
        nats_url = 'nats:4222',
        nats_subjects = 'log_subj',
        nats_format = 'JSONEachRow',
        nats_schema = 'level,text';
