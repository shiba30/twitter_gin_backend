-- name: InsertReply :one
INSERT INTO tweets (
    user_id,
    content,
    image_path,
    reply_to
) VALUES (
    $1, $2, $3, $4
) RETURNING user_id, content, image_path;