package flagsParser

import (
	"flag"
)

const DefaultPredictionLength = 60

type Flags struct {
	Model            string
	Source           string
	AggregateBy      string
	PredictionLength int64
}

func ParseFlags() *Flags {
	model := flag.String("model", "linearExtrapolation", "Model to use for prediction(linearExtrapolation)")
	source := flag.String("source", "", "Path to the source file")
	aggregateBy := flag.String("aggregate", "country", "Field to aggregate by(country|campaign)")
	predictionLength := flag.Int64("predictionLength", DefaultPredictionLength, "Length of prediction in days")
	flag.Parse()
	flags := Flags{
		Model:            *model,
		Source:           *source,
		AggregateBy:      *aggregateBy,
		PredictionLength: *predictionLength,
	}
	return &flags
}
