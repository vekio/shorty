-- name: CreateAPIKey :exec
INSERT INTO api_keys (id, workspace_id, name, token_hash, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: AuthenticateAPIKey :one
UPDATE api_keys
SET last_used_at = ?
WHERE token_hash = ? AND revoked_at IS NULL
RETURNING id, workspace_id, name, token_hash, created_at, last_used_at, revoked_at;

-- name: ListAPIKeys :many
SELECT id, workspace_id, name, token_hash, created_at, last_used_at, revoked_at
FROM api_keys
WHERE workspace_id = ?
ORDER BY created_at DESC, id DESC;

-- name: RevokeAPIKey :execrows
UPDATE api_keys
SET revoked_at = COALESCE(revoked_at, ?)
WHERE workspace_id = ? AND id = ?;
