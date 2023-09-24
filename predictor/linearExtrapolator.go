package predictor

import (
	"fmt"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
)

type LinearExtrapolator struct{}

func (le LinearExtrapolator) Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (PredictedLTVs, error) {
	var result = make(map[string]decimal.Decimal, len(al))
	for k, v := range al {
		predictedLTV, err := linearExtrapolation(v, decimal.NewFromInt(predictionLength))
		if err != nil {
			return nil, fmt.Errorf(ErrPredictorError.Error(), err)
		}
		result[k] = *predictedLTV
	}
	return result, nil
}

func linearExtrapolation(data []decimal.Decimal, predictLength decimal.Decimal) (*decimal.Decimal, error) {
	if predictLength.LessThan(decimal.NewFromInt(3)) {
		return nil, ErrPredictLengthTooShort
	}
	// decrease the length by one to take into account index starting from 0
	predictLength = predictLength.Sub(decimal.NewFromInt(1))

	// at least two points are needed to extrapolate
	if len(data) < 2 {
		return nil, ErrNotEnoughData
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
