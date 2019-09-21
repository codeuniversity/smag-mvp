CREATE TABLE IF NOT EXISTS follows(
   id serial PRIMARY KEY,
   from_id INTEGER REFERENCES users(id),
   to_id INTEGER REFERENCES users(id)
);

CREATE UNIQUE INDEX follows_uniq_relationship_index ON follows (from_id, to_id);
