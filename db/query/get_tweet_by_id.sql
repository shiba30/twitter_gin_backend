-- name: GetTweetByID :one
SELECT
    tweets.id,
    tweets.user_id,
    tweets.content,
    tweets.image_path,
    tweets.created_at,
    users.display_name
FROM
    tweets
JOIN
    users ON tweets.user_id = users.id
WHERE
    tweets.id = $1;
