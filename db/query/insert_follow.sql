-- name: CreateFollow :exec
INSERT INTO follows (
    follower_id,
    followee_id
) VALUES (
    $1, $2
);
