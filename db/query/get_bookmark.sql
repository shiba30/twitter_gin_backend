-- name: GetBookmark :one
SELECT
    id,
    tweet_id,
    user_id,
    created_at
FROM
    bookmarks
WHERE
    tweet_id = $1
    AND
    user_id = $2;
