-- name: DeleteRepliesByUserID :exec
DELETE FROM replies WHERE user_id = $1;
-- name: DeleteLikesByUserID :exec
DELETE FROM likes WHERE user_id = $1;
-- name: DeleteRetweetsByUserID :exec
DELETE FROM retweets WHERE user_id = $1;
-- name: DeleteBookmarksByTweetUserID :exec
DELETE FROM bookmarks WHERE tweet_id IN (SELECT id FROM tweets WHERE tweets.user_id = $1);
-- name: DeleteBookmarksByUserID :exec
DELETE FROM bookmarks WHERE user_id = $1;
-- name: DeleteTweetsByUserID :exec
DELETE FROM tweets WHERE user_id = $1;
-- name: DeleteFollowsByUserID :exec
DELETE FROM follows WHERE follower_id = $1 OR followee_id = $1;
-- name: DeleteMessagesByUserID :exec
DELETE FROM messages WHERE sender_id = $1 OR receiver_id = $1;
-- name: DeleteNotificationsByUserID :exec
DELETE FROM notifications WHERE user_id = $1;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;