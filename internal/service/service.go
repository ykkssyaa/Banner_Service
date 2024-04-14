package service

import (
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	"BannerService/pkg/logger"
)

type Services struct {
	Banner
}

func NewService(gateways *gateway.Gateways, logger *logger.Logger) *Services {
	return &Services{
		Banner: NewBannerService(gateways.Banner, gateways.BannerCache, logger),
	}
}

type Banner interface {
	CreateBanner(banner models.Banner) (int, error)
	GetBanner(tagId, featureId, limit, offset int32) ([]models.Banner, error)
	GetUserBanner(tagId, featureId int32, role string, useLastRevision bool) (models.Banner, error)
	DeleteBanner(id int32) error
	PatchBanner(banner models.Banner) error
	GetBannerVersions(id int32) ([]models.Banner, error)
	SetBannerVersion(id, version int32) error
}
