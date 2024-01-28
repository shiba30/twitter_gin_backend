-- name: DeleteFollow :exec
DELETE FROM
    follows
WHERE
    follower_id = $1
    AND followee_id = $2;