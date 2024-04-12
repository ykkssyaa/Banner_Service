package gateway

import (
	"BannerService/internal/consts"
	"BannerService/internal/models"
	"context"
	"github.com/redis/go-redis/v9"
)

type BannerRedis struct {
	cl *redis.Client
}

func NewBannerRedis(cl *redis.Client) *BannerRedis {
	return &BannerRedis{cl: cl}
}

func (b BannerRedis) Get(tagId, featureId int32) (models.Banner, error) {

	ctx := context.Background()
	var banner models.Banner

	err := b.cl.Get(ctx, GenKey(tagId, featureId)).Err()
	if err != nil {
		return models.Banner{}, err
	}

	return banner, nil
}

func (b BannerRedis) Set(banner models.Banner) error {

	ctx := context.Background()
	for _, tag := range banner.TagIds {
		err := b.cl.Set(ctx, GenKey(tag, banner.FeatureId), banner, consts.RedisBannerTTL).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
