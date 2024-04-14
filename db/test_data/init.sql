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

INSERT INTO feature VALUES (3), (2);
INSERT INTO tag VALUES (1);
INSERT INTO banner (feature_id,content,is_active,created_at,updated_at) VALUES (2, '{
  "title": "Active Banner",
  "text": "some_text",
  "url": "some_url"
}',true, '2024-04-09 19:27:46.9911275 +0000', '2024-04-09 19:27:46.9911275 +0000');
INSERT INTO banner (feature_id,content,is_active,created_at,updated_at) VALUES (3, '{
  "title": "Non active banner",
  "text": "some_text",
  "url": "some_url"
}',false, '2024-04-09 19:27:46.9911275 +0000', '2024-04-09 19:27:46.9911275 +0000');
INSERT INTO banner_tag_link (banner_id, tag_id) VALUES (1, 1), (2, 1);
INSERT INTO users VALUES (1, 1), (2, 2);
INSERT INTO roles VALUES (1, 1, true), (2, 2, false);