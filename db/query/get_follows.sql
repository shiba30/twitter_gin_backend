-- name: GetFollows :many
SELECT
    follows.follower_id,
    follows.followee_id,
    follows.created_at,
    users.display_name,
    users.profile_image
FROM
    follows
INNER JOIN
    users ON follows.followee_id = users.id
WHERE
    follower_id = $1
ORDER BY
    follows.created_at DESC;