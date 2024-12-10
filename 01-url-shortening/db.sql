CREATE TABLE short_urls (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT,
    short_url TEXT NOT NULL,
    user_id TEXT,
    created_by TEXT NOT NULL,
    created_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    last_modified_by TEXT,
    last_modified_date TIMESTAMP WITHOUT TIME ZONE
);

-- Create an index on short_url for faster lookups.
CREATE UNIQUE INDEX idx_short_urls_short_url ON short_urls (short_url);
