-- name: DeleteTweet :exec
DELETE FROM
    tweets
WHERE
    id = $1
    AND
    user_id = $2;
