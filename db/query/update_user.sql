-- name: UpdateUser :one
UPDATE
    users
SET
    activation_token = $2
    , is_active = $3
WHERE
    id = $1 RETURNING id,
    activation_token,
    is_active;