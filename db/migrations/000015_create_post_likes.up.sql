CREATE TABLE IF NOT EXISTS post_likes
(
    id         serial PRIMARY KEY,
    like_id VARCHAR UNIQUE NOT NULL,
    user_id    INTEGER REFERENCES users (id),
    post_id INTEGER REFERENCES posts(id)
);
