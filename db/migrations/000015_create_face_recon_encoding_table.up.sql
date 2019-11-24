CREATE TABLE IF NOT EXISTS faces
(
    id        serial PRIMARY KEY,
    post_id   VARCHAR UNIQUE NOT NULL,
	x         INTEGER,
	y         INTEGER,
	width     INTEGER,
	height    INTEGER,
	face_encoding JSONB
);
