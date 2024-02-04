-- name: InsertMessage :one
INSERT INTO messages (
    sender_id,
    receiver_id,
    content
) VALUES (
    $1, $2, $3
) RETURNING sender_id, receiver_id, content, created_at;