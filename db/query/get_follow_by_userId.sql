-- name: GetFollowByUserId :one
SELECT
    COUNT(*)
FROM
    follows
WHERE
    follower_id = $1
    AND
    followee_id = $2;