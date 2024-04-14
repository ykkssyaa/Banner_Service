ALTER TABLE tags_banners
ADD CONSTRAINT pk_tags_banners_all PRIMARY KEY(tag_id, banner_version, banner_id);