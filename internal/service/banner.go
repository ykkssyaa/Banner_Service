package service

import "BannerService/internal/gateway"

type BannerService struct {
	repo gateway.Banner
}

func NewBannerService(repo gateway.Banner) *BannerService {
	return &BannerService{repo: repo}
}
