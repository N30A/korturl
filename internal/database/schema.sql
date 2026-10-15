CREATE TYPE user_role AS enum('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username text UNIQUE NOT NULL CONSTRAINT chk_users_username_length,
    role user_role NOT NULL DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS api_keys (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name text NOT NULL CONSTRAINT chk_api_keys_name_length,
    key_hash text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz,
    revoked_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);

CREATE TABLE IF NOT EXISTS urls (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code text UNIQUE NOT NULL CONSTRAINT chk_urls_code_length,
    redirect_url text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS visits (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    url_id bigint NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    visited_at timestamptz NOT NULL DEFAULT now()
    -- senare: browser, os, ip, country, etc.
);

CREATE INDEX IF NOT EXISTS idx_visits_url_id_visited_at ON visits(url_id, visited_at);
