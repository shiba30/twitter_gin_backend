CREATE TABLE IF NOT EXISTS tweets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    content VARCHAR(140) NOT NULL,
    image_path VARCHAR(255),
    reply_to BIGINT,
    quote_tweet_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reply_to) REFERENCES tweets(id),
    FOREIGN KEY (quote_tweet_id) REFERENCES tweets(id)
);
