package predictor

import (
	"errors"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
)

const (
	ErrNotEnoughData = "not enough data to extrapolate"
)

type LinearExtrapolator struct{}

func (le LinearExtrapolator) Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (PredictedLTVs, error) {
	var result PredictedLTVs = make(map[string]decimal.Decimal, len(al))
	for k, v := range al {
		predictedLTV, err := linearExtrapolation(v, decimal.NewFromInt(predictionLength))
		if err != nil {
			return nil, err
		}
		result[k] = *predictedLTV
	}
	return result, nil
}

func linearExtrapolation(data []decimal.Decimal, predictLength decimal.Decimal) (*decimal.Decimal, error) {
	if len(data) < 2 {
		return nil, errors.New(ErrNotEnoughData)
	}
	x1, y1 := decimal.NewFromInt(int64(len(data)-1)), data[len(data)-1]
	x2, y2 := decimal.NewFromInt(int64(len(data)-1)), data[len(data)-1]

	// find the last two points that are not equal
	for i := len(data) - 2; i >= 0; i-- {
		if data[i].Equal(y1) {
			x1 = decimal.NewFromInt(int64(i))
			continue
		} else {
			y2 = data[i]
			x2 = decimal.NewFromInt(int64(i))
			break
		}
	}
	//  y = y1 + ((x - x1) / (x2 - x1)) * (y2 - y1)
	subX1FromX := predictLength.Sub(x1)
	subX1FromX2 := x2.Sub(x1)
	divSubs := subX1FromX.Div(subX1FromX2)
	subY1FromY2 := y2.Sub(y1)
	mulDivSubs := divSubs.Mul(subY1FromY2)
	predictedY := y1.Add(mulDivSubs)
	return &predictedY, nil
}
