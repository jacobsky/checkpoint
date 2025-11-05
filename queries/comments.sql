-- name: GetRecentComments :many
SELECT
    *
FROM
    comments
ORDER BY
    DATEDESC
LIMIT ?;
-- name: AddComment :exec
INSERT INTO comments(
    postdate,
    pinned,
    poster,
    message
)VALUES( ?, FALSE, ?, ?);
-- name: AddPinnedComment :exec
INSERT INTO comments(
    postdate,
    pinned,
    poster,
    message
)VALUES( ?, TRUE, ?, ?);
