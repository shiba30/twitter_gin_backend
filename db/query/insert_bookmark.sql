-- name: CreateBookmark :exec
INSERT INTO bookmarks (
    tweet_id,
    user_id
) VALUES (
    $1, $2
);
