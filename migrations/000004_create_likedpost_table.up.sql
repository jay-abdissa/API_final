-- Filename: migrations/000004_create_likedpost_table.up.sql

CREATE TABLE IF NOT EXISTS likedpost (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    users_id BIGINT REFERENCES users (id),
    posts_id BIGINT REFERENCES posts (id),
    UNIQUE(users_id, posts_id)
);