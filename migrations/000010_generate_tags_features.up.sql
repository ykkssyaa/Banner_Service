INSERT INTO features (title) (SELECT md5(random()::text) FROM generate_series(1, 1000));
INSERT INTO tags (title) (SELECT md5(random()::text) FROM generate_series(1, 1000));
