-- name: GetTweetsByUser :many
SELECT
    users.id AS user_id,
    users.display_name AS user_name,
    users.profile_image AS profile_image,
    tweets.id AS tweet_id,
    tweets.content AS tweet_content,
    tweets.image_path AS image_path,
    tweets.created_at AS tweet_date,
    tweets.is_retweet AS is_retweet,
    (SELECT COUNT(*) FROM follows WHERE follows.followee_id = users.id) AS followers_count,
    (SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count
FROM
    tweets
JOIN
    users ON tweets.user_id = users.id
WHERE
    tweets.user_id = $1
ORDER BY
    tweets.created_at DESC;
