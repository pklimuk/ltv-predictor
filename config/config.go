package config

import (
	"errors"
	"path/filepath"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/pklimuk/ltv-predictor/flagsParser"
	"github.com/pklimuk/ltv-predictor/outputPrinter"
	"github.com/pklimuk/ltv-predictor/predictor"
)

const (
	LTVDataLength = 7

	ErrUnknownModel                = "unknown model"
	ErrUnknownAggregateBy          = "unknown aggregation field"
	ErrUnsupportedFileFormat       = "source file format is not supported"
	ErrPredictionLengthNotPositive = "prediction length should be greater than 0"
	ErrPredictionLengthTooShort    = "prediction length should be greater than 7"
)

type AppConfig struct {
	Parser           fileParser.FileParser
	Aggregator       aggregator.Aggregator
	Predictor        predictor.Predictor
	PredictionLength int64
	OutputPrinter    outputPrinter.OutputPrinter
}

func CreateAppConfig(f *flagsParser.Flags) (*AppConfig, error) {
	parser, err := createParser(f)
	if err != nil {
		return nil, err
	}

	aggregator, err := createAggregator(f)
	if err != nil {
		return nil, err
	}

	predictor, err := createPredictor(f)
	if err != nil {
		return nil, err
	}

	outputPrinter := outputPrinter.ConsolePrinter{}

	err = validatePredictionLength(f.PredictionLength)
	if err != nil {
		return nil, err
	}
	return &AppConfig{
		Parser:           parser,
		Aggregator:       aggregator,
		Predictor:        predictor,
		OutputPrinter:    outputPrinter,
		PredictionLength: f.PredictionLength,
	}, nil
}

func createParser(f *flagsParser.Flags) (fileParser.FileParser, error) {
	switch filepath.Ext(f.Source) {
	case ".csv":
		return fileParser.CSVParser{Path: f.Source}, nil
	case ".json":
		return fileParser.JSONParser{Path: f.Source}, nil
	default:
		return nil, errors.New(ErrUnsupportedFileFormat)
	}
}

func createAggregator(f *flagsParser.Flags) (aggregator.Aggregator, error) {
	switch f.AggregateBy {
	case "country":
		return aggregator.ByCountryAggregator{}, nil
	case "campaign":
		return aggregator.ByCampaignAggregator{}, nil
	default:
		return nil, errors.New(ErrUnknownAggregateBy)
	}
}

func createPredictor(f *flagsParser.Flags) (predictor.Predictor, error) {
	switch f.Model {
	case "linearExtrapolation":
		return predictor.LinearExtrapolator{}, nil
	case "linearRegression":
		return predictor.LinearRegressor{}, nil
	default:
		return nil, errors.New(ErrUnknownModel)
	}
}

func validatePredictionLength(predictionLength int64) error {
	if predictionLength <= 0 {
		return errors.New(ErrPredictionLengthNotPositive)
	} else if predictionLength <= LTVDataLength {
		return errors.New(ErrPredictionLengthTooShort)
	}
	return nil
}
