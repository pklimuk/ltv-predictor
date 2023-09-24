package aggregator

import (
	"fmt"

	"github.com/pklimuk/ltv-predictor/fileParser"
)

type ByCountryAggregator struct{}

func (a ByCountryAggregator) AggregateRevenues(revenues []fileParser.Revenues) (AggregatedRevenuesByKey, error) {
	if len(revenues) == 0 {
		return nil, fmt.Errorf(ErrAggregatorError.Error(), ErrNoDataToAggregate)
	}
	var result AggregatedRevenuesByKey = make(map[string]AggregatedRevenues)
	for i := 0; i < len(revenues); i++ {
		rec := revenues[i]
		if ar, ok := result[rec.Country]; !ok {
			result[rec.Country] = AggregatedRevenues{
				Revenues:   rec.Revenues,
				UsersCount: rec.UsersCount,
			}
		} else {
			err := ar.addRevenues(rec.Revenues)
			if err != nil {
				return nil, fmt.Errorf(ErrAggregatorError.Error(), err)
			}
			ar.UsersCount += rec.UsersCount
			result[rec.Country] = ar
		}
	}
	return result, nil
}

func (a ByCountryAggregator) ConvertAggregatedByKeyRevenuesToLTVs(ar AggregatedRevenuesByKey) (AggregatedLTVsByKey, error) {
	ltvs, err := convertAggregatedByKeyRevenuesToLTVs(ar)
	if err != nil {
		return nil, fmt.Errorf(ErrAggregatorError.Error(), err)
	}
	return ltvs, nil
}
