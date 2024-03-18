-- name: UpdateUserProfile :one
UPDATE
    users
SET
    display_name = $2
    , bio = $3
    , location = $4
    , website = $5
    , birth_date = $6
    , header_image = $7
    , profile_image = $8
WHERE
    id = $1
RETURNING
    id,
    display_name,
    bio,
    location,
    website,
    birth_date,
    header_image,
    profile_image;