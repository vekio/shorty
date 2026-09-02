-- name: CreateLink :exec
INSERT INTO links (code, workspace_id, origin_url, created_at, visits)
VALUES (?, ?, ?, ?, ?);

-- name: GetLinkByCode :one
SELECT code, workspace_id, origin_url, created_at, visits
FROM links
WHERE workspace_id = ? AND code = ?;

-- name: ListLinks :many
SELECT code, workspace_id, origin_url, created_at, visits
FROM links
WHERE workspace_id = ?
ORDER BY created_at DESC, code DESC
LIMIT ? OFFSET ?;

-- name: CountLinks :one
SELECT count(*)
FROM links
WHERE workspace_id = ?;

-- name: UpdateLinkOrigin :execrows
UPDATE links
SET origin_url = ?
WHERE workspace_id = ? AND code = ?;

-- name: ResolveLink :one
UPDATE links
SET visits = visits + 1
WHERE code = ?
RETURNING code, workspace_id, origin_url, created_at, visits;

-- name: DeleteLink :execrows
DELETE FROM links
WHERE workspace_id = ? AND code = ?;
