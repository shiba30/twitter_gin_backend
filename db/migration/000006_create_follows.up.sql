CREATE TABLE IF NOT EXISTS follows (
    follower_id BIGINT NOT NULL, -- フォローする側のユーザID
    followee_id BIGINT NOT NULL, -- フォローされる側のユーザID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followee_id) REFERENCES users(id)
);