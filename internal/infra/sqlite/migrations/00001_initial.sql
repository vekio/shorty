-- +goose Up
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL
);

INSERT INTO workspaces (id, name, created_at)
VALUES (
    'ws_' || lower(hex(randomblob(16))),
    'default',
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
);

CREATE TABLE links (
    code TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    origin_url TEXT NOT NULL,
    created_at TEXT NOT NULL,
    visits INTEGER NOT NULL DEFAULT 0 CHECK (visits >= 0)
);

CREATE INDEX links_workspace_created_at_idx
    ON links (workspace_id, created_at DESC, code DESC);

CREATE TABLE api_keys (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    last_used_at TEXT,
    revoked_at TEXT
);

CREATE INDEX api_keys_workspace_created_at_idx
    ON api_keys (workspace_id, created_at DESC, id DESC);

-- +goose Down
DROP TABLE api_keys;
DROP TABLE links;
DROP TABLE workspaces;
