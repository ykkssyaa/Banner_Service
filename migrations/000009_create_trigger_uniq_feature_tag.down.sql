DROP TRIGGER IF EXISTS uniq_feature_tag_of_banner ON tags_banners;
DROP FUNCTION IF EXISTS check_uniq_feature_tag();

DROP TRIGGER IF EXISTS check_unique_feature_tags_trigger_update ON tags_banners;
DROP FUNCTION IF EXISTS check_unique_feature_tags_update();
