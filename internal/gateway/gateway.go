package gateway

import (
	"BannerService/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Gateways struct {
	Banner
	BannerCache
}

func NewGateway(db *sqlx.DB, redisCl *redis.Client, CacheOn bool) *Gateways {
	return &Gateways{
		Banner:      NewBannerPostgres(db),
		BannerCache: NewBannerRedis(redisCl, CacheOn),
	}
}

type Banner interface {
	CreateBanner(banner models.Banner) (int, error)
	GetBanner(tagId, featureId, limit, offset int32, isActive *bool) ([]models.Banner, error)
	DeleteBanner(id int32) error
	GetBannerById(id int32) (models.Banner, error)
	SetActiveVersion(id, version int32, isActive bool) error
}

type BannerCache interface {
	Get(tagId, featureId int32) (models.Banner, error)
	Set(banner models.Banner) error
}
