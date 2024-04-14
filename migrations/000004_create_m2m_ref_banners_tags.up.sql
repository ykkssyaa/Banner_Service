CREATE TABLE IF NOT EXISTS tags_banners(
    banner_id INT NOT NULL,
    banner_version INT NOT NULL,
    FOREIGN KEY (banner_id, banner_version) REFERENCES Banners(id, version)
        ON DELETE CASCADE ON UPDATE CASCADE,

    tag_id INT NOT NULL
        REFERENCES Tags(id) ON DELETE RESTRICT ON UPDATE CASCADE
);