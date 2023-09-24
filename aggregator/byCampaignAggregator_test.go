package aggregator

import (
	"testing"

	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestByCampaignAggregator_AggregateRevenues(t *testing.T) {
	aggregator := ByCampaignAggregator{}
	revenues := []fileParser.Revenues{
		{CampaignID: "campaign1", Revenues: []decimal.Decimal{decimal.NewFromFloat(100), decimal.NewFromFloat(200)}, UsersCount: 10},
		{CampaignID: "campaign1", Revenues: []decimal.Decimal{decimal.NewFromFloat(150), decimal.NewFromFloat(300)}, UsersCount: 15},
	}

	expectedResult := AggregatedRevenuesByKey{
		"campaign1": {
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

func TestByCampaignAggregator_AggregateRevenues_DifferentLengthError(t *testing.T) {
	aggregator := ByCampaignAggregator{}
	revenues := []fileParser.Revenues{
		{CampaignID: "campaign1", Revenues: []decimal.Decimal{decimal.NewFromFloat(100), decimal.NewFromFloat(200)}, UsersCount: 10},
		{CampaignID: "campaign1", Revenues: []decimal.Decimal{decimal.NewFromFloat(150)}, UsersCount: 15},
	}

	result, err := aggregator.AggregateRevenues(revenues)
	assert.Error(t, err)
	assert.Equal(t, "aggregator error: ltv and revenues slices have different length", err.Error())
	assert.Nil(t, result)
}

func TestByCampaignAggregator_AggregateRevenues_No_Data(t *testing.T) {
	aggregator := ByCampaignAggregator{}
	revenues := []fileParser.Revenues{}

	result, err := aggregator.AggregateRevenues(revenues)
	assert.Error(t, err)
	assert.Equal(t, "aggregator error: no data to aggregate", err.Error())
	assert.Nil(t, result)
}

func TestByCampaignAggregator_ConvertAggregatedByKeyRevenuesToLTVs(t *testing.T) {
	aggregator := ByCampaignAggregator{}
	aggregatedRevenues := AggregatedRevenuesByKey{
		"campaign1": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(250), decimal.NewFromFloat(500)},
			UsersCount: 25,
		},
	}

	expectedResult := AggregatedLTVsByKey{
		"campaign1": []decimal.Decimal{decimal.NewFromFloat(10), decimal.NewFromFloat(20)},
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

func TestByCampaignAggregator_ConvertAggregatedByKeyRevenuesToLTVs_DivisionByZero(t *testing.T) {
	aggregator := ByCampaignAggregator{}
	aggregatedRevenues := AggregatedRevenuesByKey{
		"campaign1": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(250), decimal.NewFromFloat(500)},
			UsersCount: 0,
		},
	}

	result, err := aggregator.ConvertAggregatedByKeyRevenuesToLTVs(aggregatedRevenues)
	assert.Error(t, err)
	assert.Equal(t, "aggregator error: division by zero", err.Error())
	assert.Nil(t, result)
}
