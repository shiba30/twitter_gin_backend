-- name: GetMessages :many
SELECT
    sender_id,
    receiver_id,
    content,
    created_at
FROM
    messages
WHERE
    (sender_id = $1 AND receiver_id = $2)
    OR (sender_id = $2 AND receiver_id = $1)
ORDER BY
    created_at;