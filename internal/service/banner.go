package service

import (
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	sErr "BannerService/pkg/serverError"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"net/http"
)

type BannerService struct {
	repo gateway.Banner
}

func NewBannerService(repo gateway.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func (p *BannerService) CreateBanner(banner models.Banner) (int, error) {

	banner.Version = 1
	id, err := p.repo.CreateBanner(banner)
	if err != nil {

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// Handling violation of sql constraint (unique)
			if pqErr.Code == "23505" {
				return 0, sErr.ServerError{
					Message:    "tags_ids has duplicates",
					StatusCode: http.StatusBadRequest,
				}
				// Handling violation of sql constraint (foreign key)
			} else if pqErr.Code == "23503" {
				return 0, sErr.ServerError{
					Message:    "Reference to a non-existent object(tag or feature)",
					StatusCode: http.StatusBadRequest,
				}
				// Handling violation of sql constraint (uniq feature and tag trigger)
			} else if pqErr.Message == "Tag already exists for this banner and feature." {
				return 0, sErr.ServerError{
					Message:    "Tag already exists for this banner and feature",
					StatusCode: http.StatusBadRequest,
				}
			}
		} else {
			return 0, sErr.ServerError{
				Message:    "Error with creating banner",
				StatusCode: http.StatusInternalServerError,
			}
		}

	}

	return id, nil
}

func (p *BannerService) GetBanner(tagId, featureId, limit, offset int32) ([]models.Banner, error) {

	if tagId < 0 {
		tagId = 0
	}
	if featureId < 0 {
		featureId = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}

	res, err := p.repo.GetBanner(tagId, featureId, limit, offset)
	if err != nil {
		return nil, sErr.ServerError{
			Message:    "Error with getting banner",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return res, nil

}

func (p *BannerService) DeleteBanner(id int32) error {

	if id <= 0 {
		return sErr.ServerError{
			Message:    "Bad Request: wrong id value",
			StatusCode: http.StatusBadRequest,
		}
	}

	err := p.repo.DeleteBanner(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sErr.ServerError{
				Message:    "",
				StatusCode: http.StatusNotFound,
			}
		} else {
			return sErr.ServerError{
				Message:    "Error with deleting banner",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	return nil
}

func (p *BannerService) PatchBanner(banner models.Banner) error {

	if banner.Id <= 0 {
		return sErr.ServerError{
			Message:    "Bad Request: wrong id value",
			StatusCode: http.StatusBadRequest,
		}
	}

	if banner.FeatureId == 0 && len(banner.TagIds) == 0 && len(banner.Content) == 0 {
		return sErr.ServerError{
			Message:    "Bad Request: nothing to update",
			StatusCode: http.StatusBadRequest,
		}
	}

	oldBanner, err := p.repo.GetBannerById(banner.Id)
	if err != nil {
		return sErr.ServerError{
			Message:    "Error with getting banner",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if oldBanner.Id == 0 {
		return sErr.ServerError{
			Message:    "",
			StatusCode: http.StatusNotFound,
		}
	}

	if oldBanner.FeatureId != banner.FeatureId && banner.FeatureId != 0 {
		oldBanner.FeatureId = banner.FeatureId
	}

	if !oldBanner.TagIds.Equal(banner.TagIds) && len(banner.TagIds) != 0 {
		oldBanner.TagIds = make(models.Tags, len(banner.TagIds))
		copy(oldBanner.TagIds, banner.TagIds)
	}

	if !oldBanner.Content.Equal(banner.Content) && len(banner.Content) != 0 {
		oldBanner.Content = banner.Content
	}

	oldBanner.Version = oldBanner.Version + 1

	_, err = p.repo.CreateBanner(oldBanner)
	if err != nil {
		return sErr.ServerError{
			Message:    "Error with updating ",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
