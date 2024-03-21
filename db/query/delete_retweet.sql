-- name: DeleteRetweet :exec
DELETE FROM
    retweets
WHERE
    original_tweet_id = $1
    AND user_id = $2;

-- name: DeleteTweet :exec
DELETE FROM
    tweets
WHERE
    id = $1
    AND
    user_id = $2;
