-- name: CreateRetweet :exec
INSERT INTO retweets (
    tweet_id,
    user_id
) VALUES (
    $1, $2
);
