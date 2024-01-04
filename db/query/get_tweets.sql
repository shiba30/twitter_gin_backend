-- name: GetTweets :many
SELECT 
    users.display_name AS user_name,
    users.profile_image AS user_image,
    tweets.id AS tweet_id,
    tweets.content AS tweet_content,
    tweets.image_path AS image_path,
    tweets.created_at AS tweet_date
FROM 
    tweets
JOIN 
    users ON tweets.user_id = users.id
ORDER BY
    tweets.created_at DESC
LIMIT
    $1 OFFSET $2;
