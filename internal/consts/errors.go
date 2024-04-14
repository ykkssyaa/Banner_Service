package consts

const (
	ErrorTagIdsDuplicates             = "tags_ids has duplicates"
	ErrorNonExistentObject            = "reference to a non-existent object(tag or feature)"
	ErrorBannerWithTagAndFeatureExist = "Tag already exists for this banner and feature."

	ErrorCreatingBanner         = "Internal server error with creating banner"
	ErrorGettingBanner          = "Internal server error with getting banner"
	ErrorDeletingBanner         = "Internal server error with deleting banner"
	ErrorUpdatingBanner         = "Internal server error with updating banner"
	ErrorUpdatingStatus         = "Internal server error with updating is_active status"
	ErrrorUpdatingActiveVersion = "Internal server error with updating active version"

	ErrorHasNoBanner    = "Bad Request: there is not banner with this id and version"
	ErrorWrongId        = "Bad Request: wrong id value"
	ErrorWrongTagId     = "Bad Request: wrong tag_id value"
	ErrorWrongFeatureId = "Bad Request: wrong feature_id value"
	ErrorWrongVersion   = "Bad Request: wrong version value"

	ErrorNoRowsAffected = "error: No rows affected"
)
