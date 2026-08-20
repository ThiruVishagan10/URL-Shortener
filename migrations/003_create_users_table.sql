CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    google_subject TEXT NOT NULL UNIQUE,

    email TEXT NOT NULL,
    name TEXT NOT NULL,

    avatar_url TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);