-- name: InsertReply :one
INSERT INTO replies (
    tweet_id,
    user_id,
    content,
    image_path
) VALUES (
    $1, $2, $3, $4
) RETURNING tweet_id, user_id, content, image_path;