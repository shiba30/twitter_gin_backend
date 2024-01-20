-- name: InsertTweet :one
INSERT INTO tweets (
    user_id,
    content,
    image_path,
    is_retweet
) VALUES (
    $1, $2, $3, $4
) RETURNING id, user_id, content, image_path, is_retweet;