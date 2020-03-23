CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   user_name VARCHAR UNIQUE NOT NULL,
   real_name VARCHAR,
   avatar_url VARCHAR,
   bio text,
   crawl_ts integer
);
