-- name: GetTweets :many
SELECT
    users.id AS user_id,
    users.display_name AS user_name,
    users.profile_image AS profile_image,
    tweets.id AS tweet_id,
    tweets.content AS tweet_content,
    tweets.image_path AS image_path,
    tweets.created_at AS tweet_date,
    tweets.is_retweet AS is_retweet,
    COALESCE(retweeters.display_name, '') AS retweeter_name, -- リツイートしたユーザーの名前
    COUNT(DISTINCT replies.id) AS replies_count,
    COUNT(DISTINCT likes.id) AS likes_count,
    COUNT(DISTINCT retweets.id) AS retweets_count,
    CASE WHEN bookmarks.tweet_id IS NOT NULL THEN true ELSE false END AS is_bookmarked,
    CASE WHEN follows.follower_id IS NOT NULL THEN true ELSE false END AS is_following
FROM
    tweets
JOIN
    users ON tweets.user_id = users.id
LEFT JOIN
    replies ON tweets.id = replies.tweet_id
LEFT JOIN
    retweets ON tweets.id = retweets.tweet_id
LEFT JOIN
    users AS retweeters ON retweets.user_id = retweeters.id AND tweets.id = retweets.tweet_id -- リツイートしたユーザー
LEFT JOIN
    likes ON tweets.id = likes.tweet_id
LEFT JOIN
    bookmarks ON tweets.id = bookmarks.tweet_id AND bookmarks.user_id = $3
LEFT JOIN
    follows ON users.id = follows.followee_id AND follows.follower_id = $3
GROUP BY
    tweets.id, users.id, retweeters.id, bookmarks.tweet_id, follows.follower_id
ORDER BY
    tweets.created_at DESC
LIMIT
    $1 OFFSET $2;
