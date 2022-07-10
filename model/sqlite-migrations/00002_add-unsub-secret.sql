-- +goose Up

CREATE TABLE comments_tmp (
  id TEXT not null primary key, 
  reply_to TEXT,
  timestamp TIMESTAMP not null,
  page TEXT not null,
  content TEXT not null,
  name TEXT not null,
  email TEXT,
  notify INTEGER not null default FALSE,
  unsubscribe_secret TEXT not null default (hex(randomblob(24))),
  FOREIGN KEY(reply_to) REFERENCES comments(id)
);

INSERT INTO comments_tmp SELECT *, (hex(randomblob(24))) FROM comments;
DROP TABLE comments;
ALTER TABLE comments_tmp RENAME TO comments;

-- +goose Down

ALTER TABLE comments
  DROP COLUMN unsubscribe_secret;