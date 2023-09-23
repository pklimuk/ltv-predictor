package predictor

import (
	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
)

const (
	ErrNotEnoughData = "not enough data to extrapolate"
)

type PredictedLTVs map[string]decimal.Decimal

type Predictor interface {
	Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (PredictedLTVs, error)
}
