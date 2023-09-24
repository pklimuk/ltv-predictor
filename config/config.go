package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/pklimuk/ltv-predictor/flagsParser"
	"github.com/pklimuk/ltv-predictor/outputPrinter"
	"github.com/pklimuk/ltv-predictor/predictor"
)

const (
	LTVDataLength = 7
)

var (
	ErrConfigError                 = errors.New("config error: %w")
	ErrUnknownModel                = errors.New("unknown model")
	ErrUnknownAggregateBy          = errors.New("unknown aggregation field")
	ErrUnsupportedFileFormat       = errors.New("source file format is not supported")
	ErrPredictionLengthNotPositive = errors.New("prediction length should be greater than 0")
	ErrPredictionLengthTooShort    = errors.New("prediction length should be greater than 7")
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
		return nil, fmt.Errorf(ErrConfigError.Error(), err)
	}

	aggregator, err := createAggregator(f)
	if err != nil {
		return nil, fmt.Errorf(ErrConfigError.Error(), err)
	}

	predictor, err := createPredictor(f)
	if err != nil {
		return nil, fmt.Errorf(ErrConfigError.Error(), err)
	}

	outputPrinter := outputPrinter.ConsolePrinter{}

	err = validatePredictionLength(f.PredictionLength)
	if err != nil {
		return nil, fmt.Errorf(ErrConfigError.Error(), err)
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
		return nil, ErrUnsupportedFileFormat
	}
}

func createAggregator(f *flagsParser.Flags) (aggregator.Aggregator, error) {
	switch f.AggregateBy {
	case "country":
		return aggregator.ByCountryAggregator{}, nil
	case "campaign":
		return aggregator.ByCampaignAggregator{}, nil
	default:
		return nil, ErrUnknownAggregateBy
	}
}

func createPredictor(f *flagsParser.Flags) (predictor.Predictor, error) {
	switch f.Model {
	case "linearExtrapolation":
		return predictor.LinearExtrapolator{}, nil
	case "linearRegression":
		return predictor.LinearRegressor{}, nil
	default:
		return nil, ErrUnknownModel
	}
}

func validatePredictionLength(predictionLength int64) error {
	if predictionLength <= 0 {
		return ErrPredictionLengthNotPositive
	} else if predictionLength <= LTVDataLength {
		return ErrPredictionLengthTooShort
	}
	return nil
}
