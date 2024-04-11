package service

import (
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	sErr "BannerService/pkg/serverError"
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
		return 0, sErr.ServerError{
			Message:    "Error with creating banner",
			StatusCode: http.StatusInternalServerError,
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
