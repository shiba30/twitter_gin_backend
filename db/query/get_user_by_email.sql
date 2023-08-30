-- name: GetUserByEmail :one
SELECT
    id,
    email,
    password,
    is_active
FROM
    users
WHERE
    email = $1;