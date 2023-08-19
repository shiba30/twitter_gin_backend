-- name: GetUserInfo :one
SELECT
    id,
    email,
    password,
    phone_number,
    display_name,
    self_introduction,
    location,
    website,
    birth_date,
    profile_image,
    avatar_image,
    is_active
FROM
    users
WHERE
    email = $1;