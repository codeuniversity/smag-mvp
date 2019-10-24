CREATE TABLE IF NOT EXISTS post_tagged_users
(
    id      serial PRIMARY KEY,
    post_id INTEGER REFERENCES posts (id),
    user_id INTEGER REFERENCES users (id)
);
