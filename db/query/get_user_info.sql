-- name: GetUserInfo :one
SELECT
    id,
    email,
    password,
    phone_number,
    display_name,
    bio,
    location,
    website,
    birth_date,
    profile_image,
    header_image,
    created_at,
    is_active
FROM
    users
WHERE
    id = $1;