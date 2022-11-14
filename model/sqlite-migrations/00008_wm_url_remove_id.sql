-- +goose Up

DROP INDEX wm_send_feed_item_url_feed_item;

CREATE TABLE wm_send_feed_item_url_tmp (
  feed_item_id TEXT NOT NULL,
  url TEXT NOT NULL,
  FOREIGN KEY (feed_item_id) REFERENCES wm_send_feed_item(id)
);

INSERT INTO wm_send_feed_item_url_tmp (feed_item_id, url)
  SELECT feed_item_id, url FROM wm_send_feed_item_url;


DROP TABLE wm_send_feed_item_url;
ALTER TABLE wm_send_feed_item_url_tmp RENAME TO wm_send_feed_item_url;

CREATE INDEX wm_send_feed_item_url_feed_item ON wm_send_feed_item_url (feed_item_id);

-- +goose Down

DROP INDEX wm_send_feed_item_url_feed_item;

CREATE TABLE wm_send_feed_item_url_tmp (
  id TEXT NOT NULL PRIMARY KEY,
  feed_item_id TEXT NOT NULL,
  url TEXT NOT NULL,
  FOREIGN KEY (feed_item_id) REFERENCES wm_send_feed_item(id)
);

INSERT INTO wm_send_feed_item_url_tmp (id, feed_item_id, url)
  SELECT hex(randomblob(16)), feed_item_id, url FROM wm_send_feed_item_url;


DROP TABLE wm_send_feed_item_url;
ALTER TABLE wm_send_feed_item_url_tmp RENAME TO wm_send_feed_item_url;

CREATE INDEX wm_send_feed_item_url_feed_item ON wm_send_feed_item_url (feed_item_id);