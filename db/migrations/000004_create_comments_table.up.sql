CREATE TABLE IF NOT EXISTS comments(
   id serial PRIMARY KEY,
   post_id INTEGER REFERENCES posts(id),
   comment_text text,
   owner_user_id INTEGER REFERENCES users(id)
);
