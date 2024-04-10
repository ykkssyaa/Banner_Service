package gateway

import (
	"github.com/jmoiron/sqlx"
)

type Gateways struct {
	BannerGateway
}

func NewGateway(db *sqlx.DB) *Gateways {
	return &Gateways{NewBannerPostgres(db)}
}

type BannerGateway interface {
}
