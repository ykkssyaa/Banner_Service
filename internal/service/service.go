package service

import (
	"BannerService/internal/gateway"
	"BannerService/internal/models"
)

type Services struct {
	Banner
}

func NewService(gateways *gateway.Gateways) *Services {
	return &Services{
		Banner: NewBannerService(gateways.Banner),
	}
}

type Banner interface {
	CreateBanner(banner models.Banner) (int, error)
	GetBanner(tagId, featureId, limit, offset int32) ([]models.Banner, error)
}
