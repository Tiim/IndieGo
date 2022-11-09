-- +goose Up

ALTER TABLE webmentions ADD COLUMN page TEXT DEFAULT "";

UPDATE webmentions SET page = substr(replace(replace(target, "https://", ""), "http://", ""), instr(replace(replace(target, "https://", ""), "http://", ""), "/") + 1);

-- +goose Down

ALTER TABLE webmentions DROP COLUMN page;

