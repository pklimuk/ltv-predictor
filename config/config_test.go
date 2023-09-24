package config

import (
	"testing"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/pklimuk/ltv-predictor/flagsParser"
	"github.com/pklimuk/ltv-predictor/outputPrinter"
	"github.com/pklimuk/ltv-predictor/predictor"
	"github.com/stretchr/testify/assert"
)

func TestCreateAppConfig(t *testing.T) {
	tests := []struct {
		name              string
		flags             *flagsParser.Flags
		expectedConfig    *AppConfig
		expectedErrString string
	}{
		{
			name: "Valid config",
			flags: &flagsParser.Flags{
				Source:           "data.csv",
				AggregateBy:      "country",
				Model:            "linearExtrapolation",
				PredictionLength: 10,
			},
			expectedConfig: &AppConfig{
				Parser:           fileParser.CSVParser{Path: "data.csv"},
				Aggregator:       aggregator.ByCountryAggregator{},
				Predictor:        predictor.LinearExtrapolator{},
				OutputPrinter:    outputPrinter.ConsolePrinter{},
				PredictionLength: 10,
			},
			expectedErrString: "",
		},
		{
			name: "Invalid prediction length",
			flags: &flagsParser.Flags{
				Source:           "data.csv",
				AggregateBy:      "country",
				Model:            "linearExtrapolation",
				PredictionLength: -1,
			},
			expectedConfig:    nil,
			expectedErrString: "config error: prediction length should be greater than 0",
		},
		{
			name: "Unknown model",
			flags: &flagsParser.Flags{
				Source:           "data.csv",
				AggregateBy:      "country",
				Model:            "invalidModel",
				PredictionLength: 10,
			},
			expectedConfig:    nil,
			expectedErrString: "config error: unknown model",
		},
		{
			name: "Invalid file format",
			flags: &flagsParser.Flags{
				Source:           "data.txt",
				AggregateBy:      "country",
				Model:            "linearExtrapolation",
				PredictionLength: 10,
			},
			expectedConfig:    nil,
			expectedErrString: "config error: source file format is not supported",
		},
		{
			name: "Invalid aggregation field",
			flags: &flagsParser.Flags{
				Source:           "data.csv",
				AggregateBy:      "invalidField",
				Model:            "linearExtrapolation",
				PredictionLength: 10,
			},
			expectedConfig:    nil,
			expectedErrString: "config error: unknown aggregation field",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := CreateAppConfig(test.flags)

			if test.expectedErrString != "" {
				assert.Error(t, err)
				assert.Nil(t, config)
				assert.EqualError(t, err, test.expectedErrString)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, test.expectedConfig.Parser, config.Parser)
				assert.Equal(t, test.expectedConfig.Aggregator, config.Aggregator)
				assert.Equal(t, test.expectedConfig.Predictor, config.Predictor)
				assert.Equal(t, test.expectedConfig.OutputPrinter, config.OutputPrinter)
				assert.Equal(t, test.expectedConfig.PredictionLength, config.PredictionLength)
			}
		})
	}
}

func TestCreateParser(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		expectedType fileParser.FileParser
		expectedErr  error
	}{
		{"CSV file", "file.csv", fileParser.CSVParser{}, nil},
		{"JSON file", "file.json", fileParser.JSONParser{}, nil},
		{"Unknown file format", "file.txt", nil, ErrUnsupportedFileFormat},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := &flagsParser.Flags{Source: test.source}
			parser, err := createParser(flags)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, parser)
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parser)
				assert.IsType(t, test.expectedType, parser)
			}
		})
	}
}

func TestCreateAggregator(t *testing.T) {
	tests := []struct {
		name         string
		aggregateBy  string
		expectedType aggregator.Aggregator
		expectedErr  error
	}{
		{"Country aggregator", "country", aggregator.ByCountryAggregator{}, nil},
		{"Campaign aggregator", "campaign", aggregator.ByCampaignAggregator{}, nil},
		{"Unknown aggregator", "unknown", nil, ErrUnknownAggregateBy},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := &flagsParser.Flags{AggregateBy: test.aggregateBy}
			aggregator, err := createAggregator(flags)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, aggregator)
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, aggregator)
				assert.IsType(t, test.expectedType, aggregator)
			}
		})
	}
}

func TestCreatePredictor(t *testing.T) {
	tests := []struct {
		name         string
		model        string
		expectedType predictor.Predictor
		expectedErr  error
	}{
		{"Linear extrapolation", "linearExtrapolation", predictor.LinearExtrapolator{}, nil},
		{"Linear regression", "linearRegression", predictor.LinearRegressor{}, nil},
		{"Unknown model", "unknown", nil, ErrUnknownModel},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := &flagsParser.Flags{Model: test.model}
			predictor, err := createPredictor(flags)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, predictor)
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, predictor)
				assert.IsType(t, test.expectedType, predictor)
			}
		})
	}
}

func TestValidatePredictionLength(t *testing.T) {
	tests := []struct {
		name             string
		predictionLength int64
		expectedErr      error
	}{
		{"Negative prediction length", -1, ErrPredictionLengthNotPositive},
		{"Zero prediction length", 0, ErrPredictionLengthNotPositive},
		{"Too short prediction length", 7, ErrPredictionLengthTooShort},
		{"Valid prediction length", 10, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validatePredictionLength(test.predictionLength)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
