-- +goose Up
CREATE TABLE IF NOT EXISTS webmentions_queue (
  id TEXT not null primary key,
  timestamp TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  tries INTEGER not null DEFAULT 0,
  next_try TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  source TEXT not null,
  target TEXT not null
);

CREATE TABLE IF NOT EXISTS webmentions (
  id TEXT not null primary key,
  ts_created TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  ts_updated TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  source TEXT not null,
  target TEXT not null,
  content TEXT not null
);

-- +goose Down

DROP TABLE webmentions;
DROP TABLE webmentions_queue;