-- name: GetLike :one
SELECT
    id,
    tweet_id,
    user_id,
    created_at
FROM
    likes
WHERE
    tweet_id = $1
    AND
    user_id = $2;
