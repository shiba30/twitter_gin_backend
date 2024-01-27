-- name: DeleteLike :exec
DELETE FROM
    likes
WHERE
    tweet_id = $1
    AND
    user_id = $2;
