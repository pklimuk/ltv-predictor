package predictor

import (
	"fmt"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat"
)

type LinearRegressor struct{}

func (lr LinearRegressor) Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (PredictedLTVs, error) {
	var result = make(map[string]decimal.Decimal, len(al))
	for k, v := range al {
		predictedLTV, err := linearRegression(v, float64(predictionLength))
		if err != nil {
			return nil, fmt.Errorf(ErrPredictorError.Error(), err)
		}
		result[k] = *predictedLTV
	}
	return result, nil
}

func linearRegression(data []decimal.Decimal, predictLength float64) (*decimal.Decimal, error) {
	if predictLength < 3 {
		return nil, ErrPredictLengthTooShort
	}
	// decrease the length by one to take into account index starting from 0
	predictLength--

	// at least two points are needed for prediction
	if len(data) < 2 {
		return nil, ErrNotEnoughData
	}

	// conversion to float64 could affect the precision, but it is not critical for this task
	ys := prepareData(data)

	xs := make([]float64, len(ys))
	for i := 0; i < len(ys); i++ {
		xs[i] = float64(i)
	}

	// y = alpha + beta*x
	alpha, beta := stat.LinearRegression(xs, ys, nil, false)
	predictedY := alpha + beta*predictLength
	predictedYDecimal := decimal.NewFromFloat(predictedY)
	return &predictedYDecimal, nil
}

// prepareData converts data to float64 and leaves only changing values
func prepareData(data []decimal.Decimal) []float64 {
	dataFloat := make([]float64, 0)
	prevValue := data[0]
	prevValueFloat, _ := prevValue.Float64()
	dataFloat = append(dataFloat, prevValueFloat)
	for i := 1; i < len(data); i++ {
		if data[i].Equal(prevValue) {
			break
		}
		valueFloat, _ := data[i].Float64()
		dataFloat = append(dataFloat, valueFloat)
		prevValue = data[i]
	}
	return dataFloat
}
