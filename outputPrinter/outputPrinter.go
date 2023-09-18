package outputPrinter

import (
	"fmt"
	"slices"

	"github.com/pklimuk/ltv-predictor/predictor"
)

type OutputPrinter interface {
	Print(data predictor.PredictedLTVs)
}

type ConsolePrinter struct{}

func (p ConsolePrinter) Print(data predictor.PredictedLTVs) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, data[k].Round(2))
	}
}
