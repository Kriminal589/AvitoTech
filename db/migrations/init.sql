ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS banner (
                                      banner_id BIGSERIAL PRIMARY KEY,
                                      feature_id BIGINT,
                                      content jsonb,
                                      is_active bool,
                                      created_at TIMESTAMP,
                                      updated_at TIMESTAMP,
                                      UNIQUE (banner_id, feature_id)
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS banner_tag_link (
    id BIGSERIAL PRIMARY KEY,
    banner_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    CONSTRAINT fk_banner_id FOREIGN KEY (banner_id)
                                           REFERENCES banner(banner_id)
                                           ON UPDATE NO ACTION
                                           ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY
) TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    admin BOOLEAN
) TABLESPACE pg_default;

INSERT INTO users VALUES (1), (2);
INSERT INTO roles VALUES (1, 1, true), (2, 2, false);