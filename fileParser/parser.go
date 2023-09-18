package fileParser

import (
	"github.com/shopspring/decimal"
)

type Revenues struct {
	Revenues   []decimal.Decimal
	Country    string
	CampaignID string
	UsersCount int64
}

type FileParser interface {
	Parse() ([]Revenues, error)
}
