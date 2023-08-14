-- name: GetUserByActivationToken :one
SELECT
    id
FROM
    users
WHERE
    activation_token = $1;