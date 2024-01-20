-- name: CreateLike :exec
INSERT INTO likes (
    tweet_id,
    user_id
) VALUES (
    $1, $2
);
