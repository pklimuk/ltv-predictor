package fileParser

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrCantOpenFile = errors.New("can't open specified file(%s)")
	ErrParsingError = errors.New("parsing error: %w")
	ErrCantReadData = errors.New("can't read data from file: %w")
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
