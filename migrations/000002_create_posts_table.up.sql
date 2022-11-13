-- Filename: migrations/000002_create_posts_table.up.sql

CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    content text NOT NULL,
    version integer NOT NULL DEFAULT 1
);