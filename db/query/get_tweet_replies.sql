-- name: GetTweetDetailReply :many
SELECT
    users.id AS user_id,
    users.display_name AS user_name,
    users.profile_image AS profile_image,
    replies.id AS reply_id,
    replies.tweet_id AS tweet_id,
    replies.content AS tweet_content,
    replies.image_path AS image_path,
    replies.created_at AS tweet_date,
    (SELECT COUNT(*) FROM likes WHERE likes.tweet_id = tweets.id) AS likes_count,
    (SELECT COUNT(*) FROM retweets WHERE retweets.tweet_id = tweets.id) AS retweets_count
FROM
    replies
JOIN
    users ON replies.user_id = users.id
JOIN
    tweets ON replies.tweet_id = tweets.id
WHERE
    replies.tweet_id = $1
ORDER BY
    replies.created_at DESC;
