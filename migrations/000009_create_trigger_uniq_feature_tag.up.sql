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