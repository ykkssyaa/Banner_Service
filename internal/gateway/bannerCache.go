package gateway

import (
	"BannerService/internal/consts"
	"BannerService/internal/models"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type BannerRedis struct {
	cl      *redis.Client
	CacheOn bool
}

func NewBannerRedis(cl *redis.Client, CacheOn bool) *BannerRedis {
	return &BannerRedis{cl: cl, CacheOn: CacheOn}
}

func (b BannerRedis) Get(tagId, featureId int32) (models.Banner, error) {

	if !b.CacheOn {
		return models.Banner{}, errors.New("error: Cache is disabled")
	}

	ctx := context.Background()
	var banner models.Banner

	err := b.cl.Get(ctx, GenKey(tagId, featureId)).Scan(&banner)
	if err != nil {
		return models.Banner{}, err
	}

	return banner, nil
}

func (b BannerRedis) Set(banner models.Banner) error {

	if !b.CacheOn {
		return errors.New("error: Cache is disabled")
	}

	ctx := context.Background()
	for _, tag := range banner.TagIds {
		err := b.cl.Set(ctx, GenKey(tag, banner.FeatureId), banner, consts.RedisBannerTTL).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
