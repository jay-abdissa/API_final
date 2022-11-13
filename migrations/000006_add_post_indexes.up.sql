-- Filename: migrations/000006_add_post_indexes.up.sql

CREATE INDEX IF NOT EXISTS forums_title_idx ON forums USING GIN(to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS forums_content_idx ON forums USING GIN(to_tsvector('simple', content));