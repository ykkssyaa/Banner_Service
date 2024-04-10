package service

import "BannerService/internal/gateway"

type Services struct {
	Banner
}

func NewService(gateways *gateway.Gateways) *Services {
	return &Services{
		Banner: NewBannerService(gateways.Banner),
	}
}

type Banner interface {
}
