CREATE TABLE short_urls (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT,
    short_url TEXT NOT NULL,
    user_id VARCHAR(255),
    qr_code TEXT,
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    last_modified_by VARCHAR(255),
    last_modified_date TIMESTAMP WITHOUT TIME ZONE
);

-- Create an index on short_url for faster lookups.
CREATE UNIQUE INDEX idx_short_urls_short_url ON short_urls (short_url);
