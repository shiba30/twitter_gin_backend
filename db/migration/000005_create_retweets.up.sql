CREATE TABLE IF NOT EXISTS retweets (
    id BIGSERIAL PRIMARY KEY,
    tweet_id BIGINT NOT NULL REFERENCES tweets(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tweet_id, user_id)
);
