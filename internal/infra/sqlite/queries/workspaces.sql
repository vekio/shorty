-- name: GetWorkspaceByName :one
SELECT id, name
FROM workspaces
WHERE name = ?;
