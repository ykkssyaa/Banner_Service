CREATE TABLE IF NOT EXISTS Tags(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) UNIQUE NOT NULL
);
