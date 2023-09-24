package main

import (
	"log"

	"github.com/pklimuk/ltv-predictor/config"
	flagsParser "github.com/pklimuk/ltv-predictor/flagsParser"
	"github.com/pklimuk/ltv-predictor/processor"
)

func main() {
	flags := flagsParser.ParseFlags()

	appConfig, err := config.CreateAppConfig(flags)
	if err != nil {
		log.Fatalf("An error occurred during configuration:\n\t%v", err)
	}

	processor := processor.Processor{
		Parser:           appConfig.Parser,
		Aggregator:       appConfig.Aggregator,
		Predictor:        appConfig.Predictor,
		PredictionLength: appConfig.PredictionLength,
		OutputPrinter:    appConfig.OutputPrinter,
	}

	err = processor.Process()
	if err != nil {
		log.Fatalf("An error occurred during processing:\n\t%v", err)
	}
}
