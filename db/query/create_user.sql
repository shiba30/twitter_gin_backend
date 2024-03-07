-- name: CreateUser :one
INSERT INTO users (
    email,
    password,
    display_name,
    birth_date
) VALUES (
    $1, $2, $3, $4
) RETURNING id, email, password, display_name, birth_date;
