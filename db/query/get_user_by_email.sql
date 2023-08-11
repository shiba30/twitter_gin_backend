-- name: GetUserByEmail :one
SELECT
    id,
    email,
    password
FROM
    users
WHERE
    email = $1;