-- name: GetTweetDetail :one
SELECT
    users.id,
    users.display_name,
    users.profile_image AS user_image,
    users.avatar_image AS avatar_image,
    tweets.id AS tweet_id,
    tweets.content AS tweet_content,
    tweets.image_path AS image_path,
    tweets.created_at AS tweet_date
FROM
    tweets
JOIN
    users ON tweets.user_id = users.id
WHERE
    tweets.id = $1;
