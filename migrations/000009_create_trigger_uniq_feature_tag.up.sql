CREATE OR REPLACE FUNCTION check_uniq_feature_tag() RETURNS trigger as $$
    BEGIN
        IF EXISTS(
            SELECT * FROM tags_banners tb
            JOIN banners b on tb.banner_id = b.id and tb.banner_version = b.version
            WHERE tb.banner_id <> NEW.banner_id AND tb.tag_id = NEW.tag_id
              AND b.feature_id = (SELECT feature_id FROM banners WHERE id = NEW.banner_id AND version = NEW.banner_version)
        ) THEN
            RAISE EXCEPTION 'Tag already exists for this banner and feature.';
        END IF;

        return NEW;

    END;
$$
language plpgsql;


CREATE OR REPLACE TRIGGER uniq_feature_tag_of_banner
BEFORE INSERT ON tags_banners
FOR EACH ROW
EXECUTE FUNCTION check_uniq_feature_tag();


CREATE OR REPLACE FUNCTION check_unique_feature_tags_update() RETURNS trigger AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM tags_banners tb
        JOIN banners b ON tb.banner_id = b.id
        WHERE tb.tag_id IN (SELECT tag_id FROM tags_banners WHERE banner_id = NEW.id AND banner_version = NEW.version)
          AND b.feature_id = NEW.feature_id
          AND tb.banner_id <> NEW.id AND tb.banner_version = NEW.version
    ) THEN
        RAISE EXCEPTION 'Another banner with the same feature and tags already exists.';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER check_unique_feature_tags_trigger_update
BEFORE UPDATE OF feature_id ON banners
FOR EACH ROW
EXECUTE FUNCTION check_unique_feature_tags_update();