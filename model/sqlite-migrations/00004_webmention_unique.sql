-- +goose Up

-- Remove duplicate webmentions
DELETE FROM webmentions WHERE EXISTS (SELECT 1 FROM webmentions w2 WHERE webmentions.source = w2.source AND webmentions.target = w2.target AND webmentions.id > w2.id);

-- Add unique constraint for source and target
CREATE UNIQUE INDEX webmentions_source_target_unique ON webmentions (source, target);

-- +goose Down

DROP INDEX webmentions_source_target_unique;
