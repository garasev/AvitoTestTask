CREATE TABLE IF NOT EXISTS slug (
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_slug (
    user_id INTEGER,
	slug_name VARCHAR(255) REFERENCES slug(name) ON DELETE CASCADE,
	dt_end TIMESTAMP DEFAULT NULL,
	PRIMARY KEY (user_id, slug_name)
);

CREATE TABLE IF NOT EXISTS archive (
	id serial PRIMARY KEY,
	user_id INTEGER,
	slug_name VARCHAR(255),
	assignment BOOL,
	dt TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avito_user (
	id serial PRIMARY KEY
);