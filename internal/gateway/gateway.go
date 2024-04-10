package gateway

import (
	"github.com/jmoiron/sqlx"
)

type Gateways struct {
	Banner
}

func NewGateway(db *sqlx.DB) *Gateways {
	return &Gateways{NewBannerPostgres(db)}
}

type Banner interface {
}
