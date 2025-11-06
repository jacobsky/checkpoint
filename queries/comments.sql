-- name: GetRecentComments :many
SELECT * FROM
    comments
ORDER BY postdate DESC
LIMIT ?
OFFSET ?;
-- name: AddComment :one
INSERT INTO comments(
    postdate,
    pinned,
    poster,
    message
)VALUES( ?, FALSE, ?, ?)
RETURNING *;
-- name: AddPinnedComment :exec
INSERT INTO comments(
    postdate,
    pinned,
    poster,
    message
)VALUES( ?, TRUE, ?, ?);
