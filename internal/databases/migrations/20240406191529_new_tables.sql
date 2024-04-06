-- +goose Up
-- +goose StatementBegin
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS tag (
    id BIGSERIAL PRIMARY KEY
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
    feature_id BIGSERIAL,
    tag_id BIGSERIAL
) TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS feature;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS banner;
-- +goose StatementEnd
