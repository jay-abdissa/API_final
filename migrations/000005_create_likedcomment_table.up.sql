-- Filename: migrations/000005_create_likedcomment_table.up.sql

CREATE TABLE IF NOT EXISTS likedcomment (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    users_id BIGINT REFERENCES users (id),
    comments_id BIGINT REFERENCES comments (id),
    UNIQUE(users_id, comments_id)
);