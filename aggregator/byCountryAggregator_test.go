package aggregator

import (
	"testing"

	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestByCountryAggregator_AggregateRevenues(t *testing.T) {
	aggregator := ByCountryAggregator{}
	revenues := []fileParser.Revenues{
		{Country: "US", Revenues: []decimal.Decimal{decimal.NewFromFloat(100), decimal.NewFromFloat(200)}, UsersCount: 10},
		{Country: "US", Revenues: []decimal.Decimal{decimal.NewFromFloat(150), decimal.NewFromFloat(300)}, UsersCount: 15},
	}

	expectedResult := AggregatedRevenuesByKey{
		"US": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(250), decimal.NewFromFloat(500)},
			UsersCount: 25,
		},
	}

	result, err := aggregator.AggregateRevenues(revenues)
	assert.NoError(t, err)
	for k, v := range result {
		for i, r := range v.Revenues {
			if r.Equal(expectedResult[k].Revenues[i]) == false {
				t.Errorf("Expected %v, got %v", expectedResult[k].Revenues[i], r)
			}
		}
	}
}

func TestByCountryAggregator_AggregateRevenues_DifferentLengthError(t *testing.T) {
	aggregator := ByCountryAggregator{}
	revenues := []fileParser.Revenues{
		{Country: "US", Revenues: []decimal.Decimal{decimal.NewFromFloat(100), decimal.NewFromFloat(200)}, UsersCount: 10},
		{Country: "US", Revenues: []decimal.Decimal{decimal.NewFromFloat(150)}, UsersCount: 15},
	}

	result, err := aggregator.AggregateRevenues(revenues)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), ErrDifferentLength)
	assert.Nil(t, result)
}

func TestByCountryAggregator_ConvertAggregatedByKeyRevenuesToLTVs(t *testing.T) {
	aggregator := ByCountryAggregator{}
	aggregatedRevenues := AggregatedRevenuesByKey{
		"US": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(250), decimal.NewFromFloat(500)},
			UsersCount: 25,
		},
	}

	expectedResult := AggregatedLTVsByKey{
		"US": []decimal.Decimal{decimal.NewFromFloat(10), decimal.NewFromFloat(20)},
	}

	result, err := aggregator.ConvertAggregatedByKeyRevenuesToLTVs(aggregatedRevenues)
	assert.NoError(t, err)
	for k, v := range result {
		for i, r := range v {
			if r.Equal(expectedResult[k][i]) == false {
				t.Errorf("Expected %v, got %v", expectedResult[k][i], r)
			}
		}
	}
}
