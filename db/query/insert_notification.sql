-- name: InsertNotification :one
INSERT INTO notifications (
    user_id,
    action_type,
    reference_id
) VALUES (
    $1, $2, $3
) RETURNING notification_id, user_id, action_type, reference_id, read, created_at;