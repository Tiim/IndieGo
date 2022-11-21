-- +goose Up

-- Empty wm_send_feed_item_url and wm_send_feed_item so all webmentions get resent.

DELETE FROM wm_send_feed_item_url;
DELETE FROM wm_send_feed_item;

-- +goose Down