package predictor

import (
	"testing"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLinearExtrapolator_Predict(t *testing.T) {
	// Create a sample input for the test
	aggregatedData := aggregator.AggregatedLTVsByKey{
		"campaign1": []decimal.Decimal{decimal.NewFromInt(1), decimal.NewFromInt(2), decimal.NewFromInt(3), decimal.NewFromInt(4), decimal.NewFromInt(5),
			decimal.NewFromInt(5), decimal.NewFromInt(5)},
		"campaign2": []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20)},
	}

	// Create a LinearExtrapolator instance
	le := LinearExtrapolator{}

	predictedLTVs, err := le.Predict(aggregatedData, 60)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert the predicted LTVs
	assert.True(t, decimal.NewFromInt(60).Equal(predictedLTVs["campaign1"]))
	assert.True(t, decimal.NewFromInt(600).Equal(predictedLTVs["campaign2"]))
}

func TestLinearExtrapolator_Predict_NotEnoughData(t *testing.T) {
	// Create a sample input for the test
	aggregatedData := aggregator.AggregatedLTVsByKey{
		"campaign1": []decimal.Decimal{decimal.NewFromInt(1)},
	}

	// Create a LinearExtrapolator instance
	le := LinearExtrapolator{}

	_, err := le.Predict(aggregatedData, 60)

	// Assert that there is no error
	assert.Error(t, err)
	assert.Equal(t, ErrNotEnoughData, err.Error())
}
