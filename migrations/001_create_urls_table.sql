CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    short_id VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);