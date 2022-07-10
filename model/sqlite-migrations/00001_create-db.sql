-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
  id TEXT not null primary key, 
  reply_to TEXT,
  timestamp TIMESTAMP not null,
  page TEXT not null,
  content TEXT not null,
  name TEXT not null,
  email TEXT,
  notify INTEGER not null default FALSE,
  FOREIGN KEY(reply_to) REFERENCES comments(id)
);

-- +goose Down
