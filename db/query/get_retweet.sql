-- name: GetRetweet :one
SELECT
    id,
    tweet_id,
    user_id,
    original_tweet_id,
    created_at
FROM
    retweets
WHERE
    original_tweet_id = $1
    AND
    user_id = $2;
