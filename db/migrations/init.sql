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
    banner_id BIGSERIAL PRIMARY KEY,
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

INSERT INTO feature VALUES (1);
INSERT INTO users VALUES (1, 1), (2, 2);
INSERT INTO roles VALUES (1, 1, true), (2, 2, false);