package gateway

import (
	"BannerService/internal/models"
	"github.com/jmoiron/sqlx"
)

type Gateways struct {
	Banner
}

func NewGateway(db *sqlx.DB) *Gateways {
	return &Gateways{NewBannerPostgres(db)}
}

type Banner interface {
	CreateBanner(banner models.Banner) (int, error)
	GetBanner(tagId, featureId, limit, offset int32) ([]models.Banner, error)
	DeleteBanner(id int32) error
}
