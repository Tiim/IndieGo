-- +goose Up

CREATE TABLE wm_send_feed_item (
  id TEXT NOT NULL PRIMARY KEY,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE wm_send_feed_item_url (
  id TEXT NOT NULL PRIMARY KEY,
  feed_item_id TEXT NOT NULL,
  url TEXT NOT NULL,
  FOREIGN KEY (feed_item_id) REFERENCES wm_send_feed_item(id)
);

CREATE INDEX wm_send_feed_item_url_feed_item ON wm_send_feed_item_url (feed_item_id);

-- +goose Down

DROP INDEX wm_send_feed_item_url_feed_item;
DROP TABLE wm_send_feed_item_url;
DROP TABLE wm_send_feed_item;
