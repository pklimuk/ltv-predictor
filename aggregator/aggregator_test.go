package aggregator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
)

func TestAggregatedRevenues_addRevenues(t *testing.T) {
	// Prepare data
	revenues := []decimal.Decimal{
		decimal.NewFromFloat(10.5499697874482206),
		decimal.NewFromFloat(21.252663605698363),
	}
	ar := &AggregatedRevenues{
		Revenues:   []decimal.Decimal{decimal.NewFromFloat(5.0), decimal.NewFromFloat(10.0)},
		UsersCount: 2,
	}

	// Call the function
	err := ar.addRevenues(revenues)

	// Assertions
	assert.Nil(t, err)
	expectedRevenues := []decimal.Decimal{
		decimal.NewFromFloat(15.5499697874482206),
		decimal.NewFromFloat(31.252663605698363),
	}
	assert.Equal(t, expectedRevenues, ar.Revenues)
}

func TestAggregatedRevenues_addRevenues_differentLength(t *testing.T) {
	// Prepare data
	revenues := []decimal.Decimal{
		decimal.NewFromFloat(10.5499697874482206),
		decimal.NewFromFloat(21.252663605698363),
	}
	ar := &AggregatedRevenues{
		Revenues:   []decimal.Decimal{decimal.NewFromFloat(5.0)},
		UsersCount: 2,
	}

	// Call the function
	err := ar.addRevenues(revenues)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, ErrDifferentLength, err)
}

func TestConvertAggregatedByKeyRevenuesToLTVs(t *testing.T) {
	// Prepare data
	ar := AggregatedRevenuesByKey{
		"key1": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(5.0), decimal.NewFromFloat(10.0)},
			UsersCount: 2,
		},
		"key2": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(15.0), decimal.NewFromFloat(30.0)},
			UsersCount: 3,
		},
	}

	// Call the function
	result, err := convertAggregatedByKeyRevenuesToLTVs(ar)

	// Assertions
	assert.Nil(t, err)
	expectedResult := AggregatedLTVsByKey{
		"key1": []decimal.Decimal{decimal.NewFromFloat(2.5), decimal.NewFromFloat(5.0)},
		"key2": []decimal.Decimal{decimal.NewFromFloat(5), decimal.NewFromFloat(10)},
	}
	for k, v := range result {
		for i := 0; i < len(v); i++ {
			if v[i].Equal(expectedResult[k][i]) == false {
				t.Errorf("Expected %v, got %v", expectedResult[k][i], v[i])
			}
		}
	}
}

func TestConvertAggregatedByKeyRevenuesToLTVs_DivisionByZero(t *testing.T) {
	// Prepare data
	ar := AggregatedRevenuesByKey{
		"key1": {
			Revenues:   []decimal.Decimal{decimal.NewFromFloat(5.0), decimal.NewFromFloat(10.0)},
			UsersCount: 0,
		},
	}

	// Call the function
	_, err := convertAggregatedByKeyRevenuesToLTVs(ar)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, ErrDivisionByZero, err)
}
