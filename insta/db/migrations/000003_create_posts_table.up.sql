CREATE TABLE IF NOT EXISTS posts
(
    id         serial PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id),
    post_id VARCHAR UNIQUE NOT NULL,
    short_code VARCHAR,
    picture_url VARCHAR,
    caption VARCHAR
);
