-- Filename: migrations/000003_create_comments_table.up.sql

CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    content text NOT NULL,
    version integer NOT NULL DEFAULT 1
);