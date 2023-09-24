package predictor

import (
	"errors"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
)

var (
	ErrPredictorError        = errors.New("predictor error: %w")
	ErrNotEnoughData         = errors.New("not enough data to make prediction")
	ErrPredictLengthTooShort = errors.New("prediction length should be greater than 2")
)

type PredictedLTVs map[string]decimal.Decimal

type Predictor interface {
	Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (PredictedLTVs, error)
}
