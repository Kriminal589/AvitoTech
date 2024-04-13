-- +goose Up
-- +goose StatementBegin
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS tag (
    id BIGSERIAL PRIMARY KEY
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS banner_tag_link (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS feature (
    id BIGSERIAL PRIMARY KEY
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tag_id BIGINT
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS banner (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGINT,
    feature_id BIGINT REFERENCES feature(id),
    content jsonb,
    is_active bool,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE (banner_id, feature_id)
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    admin BOOLEAN
) TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS feature;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS banner;
DROP TABLE IF EXISTS banner_tag_link;
-- +goose StatementEnd
