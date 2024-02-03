-- name: DeleteBookmark :exec
DELETE FROM
    bookmarks
WHERE
    tweet_id = $1
    AND
    user_id = $2;
