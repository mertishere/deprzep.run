-- name: CreateUser :one
INSERT INTO user (
	id, username, email, password, is_verified, created_at
) VALUES (
	?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM user WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM user WHERE email = ? LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM user WHERE username = ? LIMIT 1;

-- name: UpdateUserVerify :one
UPDATE user
SET
  is_verified = ?
WHERE id = ?
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE user
SET
  password = ?
WHERE id = ?
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM user WHERE id = ?;


-- name: CreateSession :one
INSERT INTO session (
	id, secret, user_id, created_at
) VALUES (
	?, ?, ?, ?
)
RETURNING *;

-- name: GetSessionByUserId :many
SELECT * FROM session WHERE user_id = ?;

-- name: GetSessionById :one
SELECT * FROM session WHERE user_id = ? AND id = ? LIMIT 1;

-- name: DeleteSessionByUserId :exec
DELETE FROM session WHERE user_id = ?;

-- name: DeleteSessionById :exec
DELETE FROM session WHERE id = ?;


-- name: CreateVerification :one
INSERT INTO verify (
	id, code, user_id, created_at
) VALUES (
	?, ?, ?, ?
)
RETURNING *;

-- name: GetVerificationByUserId :one
SELECT * FROM verify WHERE user_id = ? LIMIT 1;

-- name: GetVerificationById :one
SELECT * FROM verify WHERE user_id = ? AND id = ? LIMIT 1;

-- name: GetVerificationByCode :one
SELECT * FROM verify WHERE user_id = ? AND code = ? LIMIT 1;

-- name: DeleteVerification :exec
DELETE FROM verify WHERE user_id = ? AND code = ?;


-- name: CreateProject :one
INSERT INTO project (
	id, user_id, name, description, docs, tags, created_at, updated_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetProjectsByUserId :many
SELECT * FROM project WHERE user_id = ?;

-- name: GetProjectById :one
SELECT * FROM project WHERE user_id = ? AND id = ?;

-- name: GetProjectByName :one
SELECT * FROM project WHERE name = ?;

-- name: UpdateProjectById :one
UPDATE project
SET
  name = ?,
  description = ?,
  docs = ?,
  tags = ?,
  updated_at = ?
WHERE user_id = ? AND id = ?
RETURNING *;

-- name: DeleteProjects :exec
DELETE FROM project WHERE user_id = ?;

-- name: DeleteProjectById :exec
DELETE FROM project WHERE user_id = ? AND id = ?;

-- name: DeleteProjectByName :exec
DELETE FROM project WHERE user_id = ? AND name = ?;

-- name: CountProjects :one
SELECT COUNT(*) FROM project WHERE user_id = ?;



-- name: CreateRelease :one
INSERT INTO release (
	id, user_id, project_id, url, release, zig_version, hash, root_file, created_at, updated_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetRelease :one
SELECT * FROM release WHERE project_id = ? AND id = ? AND user_id = ? LIMIT 1;

-- name: GetReleaseById :one
SELECT * FROM release WHERE id = ? AND user_id = ? LIMIT 1;

-- name: GetReleasesByUser :many
SELECT * FROM release WHERE user_id = ?;

-- name: GetReleasesByProjectId :many
SELECT * FROM release WHERE project_id = ?;

-- name: GetReleasesByProjectIdAndRelease :one
SELECT * FROM release WHERE project_id = ? AND release = ? LIMIT 1;



-- name: UpdateReleaseById :one
UPDATE release
SET
  url = ?,
  release = ?,
  zig_version = ?,
  hash = ?,
  root_file = ?,
  updated_at = ?
WHERE user_id = ? AND id = ?
RETURNING *;



-- name: DeleteReleaseById :exec
DELETE FROM release WHERE project_id = ? AND user_id = ? AND id = ?;

-- name: DeleteRelease :exec
DELETE FROM release WHERE project_id = ? AND user_id = ? AND release = ?;
