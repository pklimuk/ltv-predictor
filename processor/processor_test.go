package processor

import (
	"errors"
	"testing"

	"github.com/pklimuk/ltv-predictor/aggregator"
	"github.com/pklimuk/ltv-predictor/fileParser"
	"github.com/pklimuk/ltv-predictor/predictor"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks for the aggregator, predictor, parser, and outputPrinter interfaces

type MockAggregator struct {
	mock.Mock
}

func (m *MockAggregator) AggregateRevenues(revenues []fileParser.Revenues) (aggregator.AggregatedRevenuesByKey, error) {
	args := m.Called(revenues)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(aggregator.AggregatedRevenuesByKey), args.Error(1)
}

func (m *MockAggregator) ConvertAggregatedByKeyRevenuesToLTVs(ar aggregator.AggregatedRevenuesByKey) (aggregator.AggregatedLTVsByKey, error) {
	args := m.Called(ar)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(aggregator.AggregatedLTVsByKey), args.Error(1)
}

type MockPredictor struct {
	mock.Mock
}

func (m *MockPredictor) Predict(al aggregator.AggregatedLTVsByKey, predictionLength int64) (predictor.PredictedLTVs, error) {
	args := m.Called(al, predictionLength)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(predictor.PredictedLTVs), args.Error(1)
}

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse() ([]fileParser.Revenues, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]fileParser.Revenues), args.Error(1)
}

type MockOutputPrinter struct {
	mock.Mock
}

func (m *MockOutputPrinter) Print(data predictor.PredictedLTVs) {
	m.Called(data)
}

func TestProcessor_Process(t *testing.T) {
	// Setup
	mockParser := new(MockParser)
	mockAggregator := new(MockAggregator)
	mockPredictor := new(MockPredictor)
	mockOutputPrinter := new(MockOutputPrinter)

	p := Processor{
		Parser:           mockParser,
		Aggregator:       mockAggregator,
		Predictor:        mockPredictor,
		PredictionLength: 7,
		OutputPrinter:    mockOutputPrinter,
	}

	// Test data
	revenues := []fileParser.Revenues{{Revenues: []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20), decimal.NewFromInt(30)}, Country: "US", CampaignID: "123", UsersCount: 2}}
	aggregatedRevenues := make(aggregator.AggregatedRevenuesByKey)
	aggregatedLTVs := make(aggregator.AggregatedLTVsByKey)
	predictions := make(predictor.PredictedLTVs)

	// Mock behavior
	mockParser.On("Parse").Return(revenues, nil)
	mockAggregator.On("AggregateRevenues", revenues).Return(aggregatedRevenues, nil)
	mockAggregator.On("ConvertAggregatedByKeyRevenuesToLTVs", aggregatedRevenues).Return(aggregatedLTVs, nil)
	mockPredictor.On("Predict", aggregatedLTVs, int64(7)).Return(predictions, nil)
	mockOutputPrinter.On("Print", predictions).Return()

	// Execute the method under test
	err := p.Process()

	// Assertions
	assert.NoError(t, err)
	mockParser.AssertExpectations(t)
	mockAggregator.AssertExpectations(t)
	mockPredictor.AssertExpectations(t)
	mockOutputPrinter.AssertExpectations(t)
}

func TestProcessor_Process_ErrorInParser(t *testing.T) {
	// Setup
	mockParser := new(MockParser)
	mockAggregator := new(MockAggregator)

	p := Processor{
		Parser:           mockParser,
		Aggregator:       mockAggregator,
		Predictor:        nil,
		PredictionLength: 7,
		OutputPrinter:    nil,
	}

	// Mock behavior - Error in parser
	mockParser.On("Parse").Return(nil, errors.New("error parsing"))

	// Execute the method under test
	err := p.Process()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "error parsing")
	mockParser.AssertExpectations(t)
	mockAggregator.AssertNotCalled(t, "AggregateRevenues")
}

func TestProcessor_Process_ErrorInAggregator(t *testing.T) {
	// Setup
	mockParser := new(MockParser)
	mockAggregator := new(MockAggregator)

	p := Processor{
		Parser:           mockParser,
		Aggregator:       mockAggregator,
		Predictor:        nil,
		PredictionLength: 7,
		OutputPrinter:    nil,
	}

	// Test data
	revenues := []fileParser.Revenues{{Revenues: []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20), decimal.NewFromInt(30)}, Country: "US", CampaignID: "123", UsersCount: 2}}

	// Mock behavior - Error in aggregator
	mockParser.On("Parse").Return(revenues, nil)
	mockAggregator.On("AggregateRevenues", revenues).Return(nil, errors.New("error aggregating"))

	// Execute the method under test
	err := p.Process()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "error aggregating")
	mockParser.AssertExpectations(t)
	mockAggregator.AssertExpectations(t)
}

func TestProcessor_Process_ErrorInConvertAggregatedByKeyRevenuesToLTVs(t *testing.T) {
	// Setup
	mockParser := new(MockParser)
	mockAggregator := new(MockAggregator)

	p := Processor{
		Parser:           mockParser,
		Aggregator:       mockAggregator,
		Predictor:        nil,
		PredictionLength: 7,
		OutputPrinter:    nil,
	}

	// Test data
	revenues := []fileParser.Revenues{{Revenues: []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20), decimal.NewFromInt(30)}, Country: "US", CampaignID: "123", UsersCount: 2}}
	aggregatedRevenues := make(aggregator.AggregatedRevenuesByKey)

	// Mock behavior - Error in aggregator
	mockParser.On("Parse").Return(revenues, nil)
	mockAggregator.On("AggregateRevenues", revenues).Return(aggregatedRevenues, nil)
	mockAggregator.On("ConvertAggregatedByKeyRevenuesToLTVs", aggregatedRevenues).Return(nil, errors.New("error converting"))

	// Execute the method under test
	err := p.Process()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "error converting")
	mockParser.AssertExpectations(t)
	mockAggregator.AssertExpectations(t)
}

func TestProcessor_Process_ErrorInPredictor(t *testing.T) {
	// Setup
	mockParser := new(MockParser)
	mockAggregator := new(MockAggregator)
	mockPredictor := new(MockPredictor)

	p := Processor{
		Parser:           mockParser,
		Aggregator:       mockAggregator,
		Predictor:        mockPredictor,
		PredictionLength: 7,
		OutputPrinter:    nil,
	}

	// Test data
	revenues := []fileParser.Revenues{{Revenues: []decimal.Decimal{decimal.NewFromInt(10), decimal.NewFromInt(20), decimal.NewFromInt(30)}, Country: "US", CampaignID: "123", UsersCount: 2}}
	aggregatedRevenues := make(aggregator.AggregatedRevenuesByKey)
	aggregatedLTVs := make(aggregator.AggregatedLTVsByKey)

	// Mock behavior - Error in aggregator
	mockParser.On("Parse").Return(revenues, nil)
	mockAggregator.On("AggregateRevenues", revenues).Return(aggregatedRevenues, nil)
	mockAggregator.On("ConvertAggregatedByKeyRevenuesToLTVs", aggregatedRevenues).Return(aggregatedLTVs, nil)
	mockPredictor.On("Predict", aggregatedLTVs, int64(7)).Return(nil, errors.New("error predicting"))

	// Execute the method under test
	err := p.Process()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "error predicting")
	mockParser.AssertExpectations(t)
	mockAggregator.AssertExpectations(t)
	mockPredictor.AssertExpectations(t)
}
