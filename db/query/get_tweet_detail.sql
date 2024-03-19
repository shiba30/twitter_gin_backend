-- name: GetTweetDetail :one
SELECT
    users.id AS user_id,
    users.display_name AS user_name,
    users.profile_image AS profile_image,
    tweets.id AS tweet_id,
    tweets.content AS tweet_content,
    tweets.image_path AS image_path,
    tweets.created_at AS tweet_date,
    tweets.is_retweet AS is_retweet,
    COUNT(DISTINCT replies.id) AS replies_count,
    COUNT(DISTINCT likes.id) AS likes_count,
    COUNT(DISTINCT retweets.id) AS retweets_count,
    CASE WHEN bookmarks.tweet_id IS NOT NULL THEN true ELSE false END AS is_bookmarked
FROM
    tweets
JOIN
    users ON tweets.user_id = users.id
LEFT JOIN
    replies ON tweets.id = replies.tweet_id
LEFT JOIN
    retweets ON tweets.id = retweets.tweet_id
LEFT JOIN
    likes ON tweets.id = likes.tweet_id
LEFT JOIN
    bookmarks ON tweets.id = bookmarks.tweet_id AND bookmarks.user_id = $2
WHERE
    tweets.id = $1
GROUP BY
    users.id, tweets.id, bookmarks.tweet_id;
