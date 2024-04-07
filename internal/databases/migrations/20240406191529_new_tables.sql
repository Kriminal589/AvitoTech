-- +goose Up
-- +goose StatementBegin
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS tag (
    id BIGSERIAL PRIMARY KEY
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS banner_tmp (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGSERIAL NOT NULL,
    tag_id BIGSERIAL NOT NULL
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS feature (
    id BIGSERIAL PRIMARY KEY
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    admin BOOLEAN DEFAULT FALSE NOT NULL,
    feature_id BIGSERIAL,
    tag_id BIGSERIAL
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS banner (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGSERIAL,
    feature_id BIGSERIAL REFERENCES feature(id),
    message jsonb,
    UNIQUE (banner_id, feature_id)
) TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS feature;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS banner;
-- +goose StatementEnd
