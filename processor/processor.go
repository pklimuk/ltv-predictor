package processor

import (
	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/pklimuk/ltv-predictor/outputPrinter"
	"github.com/pklimuk/ltv-predictor/predictor"
)

type Processor struct {
	Parser           fileParser.FileParser
	Aggregator       aggregator.Aggregator
	Predictor        predictor.Predictor
	PredictionLength int64
	OutputPrinter    outputPrinter.OutputPrinter
}

func (p *Processor) Process() error {
	data, err := p.Parser.Parse()
	if err != nil {
		return err
	}
	aggregatedRevenues, err := p.Aggregator.AggregateRevenues(data)
	if err != nil {
		return err
	}
	aggregatedLTVs, err := p.Aggregator.ConvertAggregatedByKeyRevenuesToLTVs(aggregatedRevenues)
	if err != nil {
		return err
	}
	predictions, err := p.Predictor.Predict(aggregatedLTVs, p.PredictionLength)
	if err != nil {
		return err
	}
	p.OutputPrinter.Print(predictions)
	return nil
}
