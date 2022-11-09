-- +goose Up

-- Add the authorname and content fields to the webmentions table in sqlite

ALTER TABLE webmentions ADD COLUMN author_name TEXT DEFAULT "";
ALTER TABLE webmentions ADD COLUMN content TEXT DEFAULT "";

-- +goose Down


ALTER TABLE webmentions DROP COLUMN author_name;
ALTER TABLE webmentions DROP COLUMN content;
