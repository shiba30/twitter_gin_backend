-- name: InsertTweet :one
INSERT INTO tweets (
    user_id,
    content,
    image_path
) VALUES (
    $1, $2, $3
) RETURNING user_id, content, image_path;