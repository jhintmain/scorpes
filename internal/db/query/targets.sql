-- name: ListTargets :many
SELECT *
FROM targets
WHERE deleted_at IS NULL
  AND is_active = true;

-- name: GetTargetByID :one
SELECT *
FROM targets
WHERE id = $1
  AND deleted_at IS NULL;

-- name: CreateTarget :one
INSERT INTO targets (name, url, method, interval_seconds, timeout_seconds)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: SoftDeleteTarget :exec
UPDATE targets
SET deleted_at = now(),
    is_active  = false
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateTarget :one
UPDATE targets
SET name = $2,
    url = $3,
    method = $4,
    interval_seconds = $5,
    timeout_seconds = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;