package aggregator

import (
	"errors"

	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/shopspring/decimal"
)

const (
	ErrDivisionByZero  = "division by zero"
	ErrDifferentLength = "ltv and revenues slices have different length"
)

type Aggregator interface {
	AggregateRevenues(revenues []fileParser.Revenues) (AggregatedRevenuesByKey, error)
	ConvertAggregatedByKeyRevenuesToLTVs(ar AggregatedRevenuesByKey) (AggregatedLTVsByKey, error)
}

type AggregatedRevenues struct {
	Revenues   []decimal.Decimal
	UsersCount int64
}
type AggregatedRevenuesByKey map[string]AggregatedRevenues
type AggregatedLTVs []decimal.Decimal
type AggregatedLTVsByKey map[string]AggregatedLTVs

func (ar *AggregatedRevenues) addUserLtvToRevenues(ltv []decimal.Decimal) error {
	if len(ltv) != len(ar.Revenues) {
		return errors.New(ErrDifferentLength)
	}
	for i := 0; i < len(ltv); i++ {
		ar.Revenues[i] = ar.Revenues[i].Add(ltv[i])
	}
	return nil
}

func convertAggregatedByKeyRevenuesToLTVs(ar AggregatedRevenuesByKey) (AggregatedLTVsByKey, error) {
	var result AggregatedLTVsByKey = make(map[string]AggregatedLTVs)
	for k, v := range ar {
		ltvs := make([]decimal.Decimal, len(v.Revenues))
		for i := 0; i < len(v.Revenues); i++ {
			if v.UsersCount == 0 {
				return nil, errors.New(ErrDivisionByZero)
			}
			ltvs[i] = v.Revenues[i].Div(decimal.NewFromInt(v.UsersCount))
		}
		result[k] = ltvs
	}
	return result, nil
}
