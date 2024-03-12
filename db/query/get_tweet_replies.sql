-- name: GetTweetDetailReply :many
SELECT
    users.id,
    users.display_name,
    users.profile_image AS profile_image,
    users.avatar_image AS avatar_image,
    replies.id AS reply_id,
    replies.tweet_id AS tweet_id,
    replies.content AS reply_content,
    replies.image_path AS image_path,
    replies.created_at AS reply_date
FROM
    replies
JOIN
    users ON replies.user_id = users.id
WHERE
    replies.tweet_id = $1
ORDER BY
    replies.created_at DESC;
