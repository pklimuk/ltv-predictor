package aggregator

import (
	"errors"

	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/shopspring/decimal"
)

var (
	ErrAggregatorError   = errors.New("aggregator error: %w")
	ErrNoDataToAggregate = errors.New("no data to aggregate")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrDifferentLength   = errors.New("ltv and revenues slices have different length")
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

func (ar *AggregatedRevenues) addRevenues(revenues []decimal.Decimal) error {
	if len(revenues) != len(ar.Revenues) {
		return ErrDifferentLength
	}
	for i := 0; i < len(revenues); i++ {
		ar.Revenues[i] = ar.Revenues[i].Add(revenues[i])
	}
	return nil
}

func convertAggregatedByKeyRevenuesToLTVs(ar AggregatedRevenuesByKey) (AggregatedLTVsByKey, error) {
	var result AggregatedLTVsByKey = make(map[string]AggregatedLTVs)
	for k, v := range ar {
		ltvs := make([]decimal.Decimal, len(v.Revenues))
		revsLen := len(v.Revenues)
		for i := 0; i < revsLen; i++ {
			if v.UsersCount == 0 {
				return nil, ErrDivisionByZero
			}
			ltvs[i] = v.Revenues[i].Div(decimal.NewFromInt(v.UsersCount))
		}
		result[k] = ltvs
	}
	return result, nil
}
