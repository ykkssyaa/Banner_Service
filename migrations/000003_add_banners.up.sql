CREATE TABLE IF NOT EXISTS Banners(
    id SERIAL,
    version INT,
    PRIMARY KEY(id, version),
    feature INT NOT NULL REFERENCES Features(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    content JSON NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at timestamp default now(),
    update_at timestamp
);

