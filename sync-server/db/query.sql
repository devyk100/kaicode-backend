-- name: UpdateSessionContent :exec
UPDATE "sessions"
SET content = $2
WHERE id = $1;