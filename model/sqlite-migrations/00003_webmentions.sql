-- +goose Up
CREATE TABLE IF NOT EXISTS webmentions_queue (
  id TEXT not null primary key,
  timestamp TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  source TEXT not null,
  target TEXT not null
);

CREATE TABLE IF NOT EXISTS webmentions_rejected (
  id TEXT not null primary key,
  timestamp TIMESTAMP not null,
  source TEXT not null,
  target TEXT not null,
  reason TEXT not null
);


CREATE TABLE IF NOT EXISTS webmentions (
  id TEXT not null primary key,
  ts_created TIMESTAMP not null,
  ts_updated TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
  source TEXT not null,
  target TEXT not null,
  deleted BOOLEAN not null DEFAULT false
);

CREATE TABLE IF NOT EXISTS domain_deny_list (
  domain TEXT not null primary key
);

-- +goose Down

DROP TABLE webmentions;
DROP TABLE webmentions_queue;
DROP TABLE webmentions_rejected;