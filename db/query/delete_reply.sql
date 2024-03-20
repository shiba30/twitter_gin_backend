-- name: DeleteReply :exec
DELETE FROM
    replies
WHERE
    id = $1;