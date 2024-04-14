package tests

import "time"
import "BannerService/internal/models"

func boolPtr(b bool) *bool {
	return &b
}

var banners = []models.Banner{
	{
		Id:        1,
		Version:   1,
		TagIds:    models.Tags{1, 2, 3},
		FeatureId: 1,
		Content:   models.ModelMap{"title": "Banner 1", "image_url": "http://example.com/banner1.jpg"},
		IsActive:  boolPtr(true),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        2,
		Version:   1,
		TagIds:    models.Tags{2, 3, 4},
		FeatureId: 2,
		Content:   models.ModelMap{"title": "Banner 2", "image_url": "http://example.com/banner2.jpg"},
		IsActive:  boolPtr(true),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
