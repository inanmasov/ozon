CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    original_url TEXT UNIQUE NOT NULL,
    short_url TEXT UNIQUE NOT NULL
);
