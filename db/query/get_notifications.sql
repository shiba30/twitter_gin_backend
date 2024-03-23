-- name: GetNotifications :many
SELECT
    n.notification_id,
    n.action_type,
    n.created_at AS notification_created_at,
    n.read,
    u.display_name AS action_user_display_name,
    u.profile_image AS action_user_profile_image,
    t.content AS related_tweet_content,
    t.image_path AS related_tweet_image_path,
    r.content AS related_reply_content
FROM
    notifications n
LEFT JOIN follows f ON n.action_type = 'follow' AND n.reference_id = f.followee_id
LEFT JOIN likes l ON n.action_type = 'like' AND n.reference_id = l.tweet_id AND n.user_id != l.user_id
LEFT JOIN replies r ON n.action_type = 'comment' AND n.reference_id = r.tweet_id AND n.user_id != r.user_id
LEFT JOIN tweets t ON t.id = l.tweet_id OR t.id = r.tweet_id
LEFT JOIN users u ON f.follower_id = u.id OR l.user_id = u.id OR r.user_id = u.id
WHERE
    n.user_id = $1
ORDER BY
    n.created_at DESC;