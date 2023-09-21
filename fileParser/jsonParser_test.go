package fileParser

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func createTempJSONFile(data string) (string, error) {
	tmpfile, err := os.CreateTemp("", "testjson_*.json")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(data)); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func removeTempJSONFile(path string) error {
	return os.Remove(path)
}

func TestJSONParser_Parse(t *testing.T) {
	// Sample JSON data
	jsonData := `[{
    "CampaignId": "9566c74d-1003-4c4d-bbbb-0407d1e2c649",
    "Country": "TR",
    "Ltv1": 1.9542502880389,
    "Ltv2": 1.9941329469784,
    "Ltv3": 3.0126373791241,
    "Ltv4": 3.1138048970185,
    "Ltv5": 3.2014612651819,
    "Ltv6": 3.7967986751124,
    "Ltv7": 4.3219611617577,
    "Users": 93
  }]`

	// Create a temporary JSON file
	tempJSONFilePath, err := createTempJSONFile(jsonData)
	if err != nil {
		t.Fatalf("Failed to create temp JSON file: %v", err)
	}
	defer removeTempJSONFile(tempJSONFilePath)

	parser := JSONParser{
		Path: tempJSONFilePath,
	}

	expectedRevenues := []Revenues{
		{
			Revenues: []decimal.Decimal{decimal.NewFromFloat(181.7452767876177), decimal.NewFromFloat(185.4543640689912), decimal.NewFromFloat(280.1752762585413),
				decimal.NewFromFloat(289.5838554227205), decimal.NewFromFloat(297.7358976619167), decimal.NewFromFloat(353.1022767854532), decimal.NewFromFloat(401.9423880434661)},
			Country:    "TR",
			CampaignID: "9566c74d-1003-4c4d-bbbb-0407d1e2c649",
			UsersCount: 93,
		},
	}

	revenues, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	assert.Equal(t, expectedRevenues, revenues)
}

func TestJSONParser_Parse_Invalid_JSON(t *testing.T) {
	t.Run("Empty JSON data", func(t *testing.T) {
		// Empty JSON data
		jsonData := `[]`

		// Create a temporary JSON file
		tempJSONFilePath, err := createTempJSONFile(jsonData)
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer removeTempJSONFile(tempJSONFilePath)

		parser := JSONParser{
			Path: tempJSONFilePath,
		}

		revenues, err := parser.Parse()
		assert.Nil(t, err)
		assert.Empty(t, revenues)
	})

	t.Run("Invalid JSON data", func(t *testing.T) {
		// Invalid JSON data
		jsonData := `[{INVALID JSON}]`

		// Create a temporary JSON file
		tempJSONFilePath, err := createTempJSONFile(jsonData)
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer removeTempJSONFile(tempJSONFilePath)

		parser := JSONParser{
			Path: tempJSONFilePath,
		}

		revenues, err := parser.Parse()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
		assert.Empty(t, revenues)
	})

	t.Run("Invalid file path", func(t *testing.T) {
		// Invalid file path
		invalidPath := "invalid/path/to/file.json"

		parser := JSONParser{
			Path: invalidPath,
		}

		revenues, err := parser.Parse()
		assert.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, revenues)
	})
}
