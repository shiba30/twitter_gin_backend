-- name: CreateRetweet :exec
INSERT INTO retweets (
    tweet_id,
    user_id,
    original_tweet_id
) VALUES (
    $1, $2, $3
);
