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
    -- eventuellt senare: browser, os, ip, country, etc.
);

CREATE INDEX IF NOT EXISTS idx_visits_url_id_visited_at ON visits(url_id, visited_at);
