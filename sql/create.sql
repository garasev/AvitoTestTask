CREATE TABLE IF NOT EXISTS slug (
	id SERIAL PRIMARY KEY,
    name VARCHAR(255) unique
);

CREATE TABLE IF NOT EXISTS user_slug (
    user_id INTEGER,
	slug_id INTEGER REFERENCES slug(id),
	dt_end TIMESTAMP DEFAULT NULL,
	PRIMARY KEY (user_id, slug_id)
);

CREATE TABLE IF NOT EXISTS archive (
	id serial PRIMARY KEY,
	user_id INTEGER,
	slug_id INTEGER,
	assignment BOOL,
	dt TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avito_user (
	id serial PRIMARY KEY
);