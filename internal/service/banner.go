package service

import (
	"BannerService/internal/consts"
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	sErr "BannerService/pkg/serverError"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"net/http"
)

type BannerService struct {
	repo  gateway.Banner
	cache gateway.BannerCache
}

func NewBannerService(repo gateway.Banner, cache gateway.BannerCache) *BannerService {
	return &BannerService{repo: repo, cache: cache}
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
					Message:    consts.ErrorTagIdsDuplicates,
					StatusCode: http.StatusBadRequest,
				}
				// Handling violation of sql constraint (foreign key)
			} else if pqErr.Code == "23503" {
				return 0, sErr.ServerError{
					Message:    consts.ErrorNonExistentObject,
					StatusCode: http.StatusBadRequest,
				}
				// Handling violation of sql constraint (uniq feature and tag trigger)
			} else if pqErr.Message == consts.ErrorBannerWithTagAndFeatureExist {
				return 0, sErr.ServerError{
					Message:    consts.ErrorBannerWithTagAndFeatureExist,
					StatusCode: http.StatusBadRequest,
				}
			}
		} else {
			return 0, sErr.ServerError{
				Message:    consts.ErrorCreatingBanner,
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

	res, err := p.repo.GetBanner(tagId, featureId, limit, offset, nil)
	if err != nil {
		return nil, sErr.ServerError{
			Message:    consts.ErrorGettingBanner,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return res, nil

}

func (p *BannerService) DeleteBanner(id int32) error {

	if id <= 0 {
		return sErr.ServerError{
			Message:    consts.ErrorWrongId,
			StatusCode: http.StatusBadRequest,
		}
	}

	err := p.repo.DeleteBanner(id)

	if err != nil {
		if err.Error() == consts.ErrorNoRowsAffected {
			return sErr.ServerError{
				Message:    "",
				StatusCode: http.StatusNotFound,
			}
		} else {
			return sErr.ServerError{
				Message:    consts.ErrorDeletingBanner,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	return nil
}

func (p *BannerService) PatchBanner(banner models.Banner) error {

	if banner.Id <= 0 {
		return sErr.ServerError{
			Message:    consts.ErrorWrongId,
			StatusCode: http.StatusBadRequest,
		}
	}

	banners, err := p.repo.GetBannersById(banner.Id, true)
	oldBanner := banners[0]
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return sErr.ServerError{
				Message:    "",
				StatusCode: http.StatusNotFound,
			}
		}

		return sErr.ServerError{
			Message:    consts.ErrorGettingBanner,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Если изменился только статус активности
	if banner.FeatureId == 0 && len(banner.TagIds) == 0 && len(banner.Content) == 0 &&
		banner.IsActive != nil && *banner.IsActive != *oldBanner.IsActive {

		err = p.repo.SetActiveVersion(oldBanner.Id, oldBanner.Version, *banner.IsActive)
	} else {

		if banner.FeatureId != 0 {
			oldBanner.FeatureId = banner.FeatureId
		}

		if len(banner.TagIds) != 0 {
			oldBanner.TagIds = make(models.Tags, len(banner.TagIds))
			copy(oldBanner.TagIds, banner.TagIds)
		}

		if len(banner.Content) != 0 {
			oldBanner.Content = banner.Content
		}

		oldBanner.Version = oldBanner.Version + 1

		_, err = p.repo.CreateBanner(oldBanner)
		if err != nil {
			return sErr.ServerError{
				Message:    consts.ErrorUpdatingBanner,
				StatusCode: http.StatusInternalServerError,
			}
		}

		var isActive bool
		if banner.IsActive == nil {
			isActive = *oldBanner.IsActive
		} else {
			isActive = *banner.IsActive
		}

		err = p.repo.SetActiveVersion(oldBanner.Id, oldBanner.Version, isActive)
	}

	if err != nil {
		return sErr.ServerError{
			Message:    consts.ErrorUpdatingStatus,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (p *BannerService) GetUserBanner(tagId, featureId int32, role string, useLastRevision bool) (models.Banner, error) {

	if tagId <= 0 {
		return models.Banner{}, sErr.ServerError{
			Message:    consts.ErrorWrongTagId,
			StatusCode: http.StatusBadRequest,
		}
	}
	if featureId <= 0 {
		return models.Banner{}, sErr.ServerError{
			Message:    consts.ErrorWrongFeatureId,
			StatusCode: http.StatusBadRequest,
		}
	}

	isActive := new(bool)
	if role != consts.AdminRole {
		*isActive = true
	} else {
		isActive = nil
	}

	var banner models.Banner
	if !useLastRevision {

		cachedBanner, err := p.cache.Get(tagId, featureId)
		if err != nil {
			// TODO: Не нарушаем работу программы, нужно логгировать об ошибке
		}
		banner = cachedBanner
	}

	if banner.Id == 0 || banner.IsActive != nil && !*banner.IsActive && role != consts.AdminRole { // Banner there are not in cache
		banners, err := p.repo.GetBanner(tagId, featureId, 1, 0, isActive)
		if err != nil {
			return models.Banner{}, sErr.ServerError{
				Message:    consts.ErrorGettingBanner,
				StatusCode: http.StatusInternalServerError,
			}
		}

		if len(banners) == 0 {
			return models.Banner{}, sErr.ServerError{
				Message:    "",
				StatusCode: http.StatusNotFound,
			}
		}

		banner = banners[0]

		if err := p.cache.Set(banner); err != nil {
			// TODO: Не нарушаем работу программы, нужно логгировать об ошибке
		}
	}

	return banner, nil
}

func (p *BannerService) GetBannerVersions(id int32) ([]models.Banner, error) {

	if id <= 0 {
		return nil, sErr.ServerError{
			Message:    consts.ErrorWrongId,
			StatusCode: http.StatusBadRequest,
		}
	}

	banners, err := p.repo.GetBannersById(id, false)
	if err != nil {
		return nil, sErr.ServerError{
			Message:    consts.ErrorGettingBanner,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if len(banners) == 0 {
		return nil, sErr.ServerError{
			Message:    "",
			StatusCode: http.StatusNotFound,
		}
	}

	return banners, nil
}

func (p *BannerService) SetBannerVersion(id, version int32) error {

	if id <= 0 {
		return sErr.ServerError{
			Message:    consts.ErrorWrongId,
			StatusCode: http.StatusBadRequest,
		}
	}

	if version <= 0 {
		return sErr.ServerError{
			Message:    consts.ErrorWrongVersion,
			StatusCode: http.StatusBadRequest,
		}
	}

	err := p.repo.SetActiveVersion(id, version, true)

	if err != nil {
		if err.Error() == consts.ErrorNoRowsAffected {
			return sErr.ServerError{
				Message:    consts.ErrorHasNoBanner,
				StatusCode: http.StatusBadRequest,
			}
		}
		return sErr.ServerError{
			Message:    consts.ErrrorUpdatingActiveVersion,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
