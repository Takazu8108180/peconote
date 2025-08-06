CREATE TABLE IF NOT EXISTS memo (
    id UUID PRIMARY KEY,
    body TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_memo_tags ON memo USING GIN (tags);
