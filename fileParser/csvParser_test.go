package fileParser

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func createTempCSVFile(data string) (string, error) {
	tmpfile, err := os.CreateTemp("", "testcsv_*.csv")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(data)); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func removeTempCSVFile(path string) error {
	return os.Remove(path)
}

func TestCSVParser_Parse(t *testing.T) {
	// Sample CSV data
	csvData := `UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7
		1,81855ad8-681d-4d86-91e9-1e00167939cb,TR,1.54996978744822,2.2526636056983,2.29863633234526,2.88400864327196,3.6960018085883,5.7144365110237,0
		2,asjdfbdd-fdg1-84gi-f9ge-aspmr9554462,IT,3.89419489191847,4.5053270142454,4.5975908308631,5.7800216081799,7.3920045214707,0,0
		3,u6jyfnon-1f1d-4d86-91e9-1e00167939cb,TR,0.87894954884928,1.01333200284909,1.03490792577263,0,0,0,0`

	// Create a temporary CSV file
	tempCSVFilePath, err := createTempCSVFile(csvData)
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer removeTempCSVFile(tempCSVFilePath)

	parser := CSVParser{
		Path: tempCSVFilePath,
	}

	expectedRevenues := []Revenues{
		{
			Revenues: []decimal.Decimal{decimal.NewFromFloat(1.54996978744822), decimal.NewFromFloat(2.2526636056983), decimal.NewFromFloat(2.29863633234526), decimal.NewFromFloat(2.88400864327196),
				decimal.NewFromFloat(3.6960018085883), decimal.NewFromFloat(5.7144365110237), decimal.NewFromFloat(5.7144365110237)},
			Country:    "TR",
			CampaignID: "81855ad8-681d-4d86-91e9-1e00167939cb",
			UsersCount: 1,
		},
		{
			Revenues: []decimal.Decimal{decimal.NewFromFloat(3.89419489191847), decimal.NewFromFloat(4.5053270142454), decimal.NewFromFloat(4.5975908308631), decimal.NewFromFloat(5.7800216081799),
				decimal.NewFromFloat(7.3920045214707), decimal.NewFromFloat(7.3920045214707), decimal.NewFromFloat(7.3920045214707)},
			Country:    "IT",
			CampaignID: "asjdfbdd-fdg1-84gi-f9ge-aspmr9554462",
			UsersCount: 1,
		},
		{
			Revenues: []decimal.Decimal{decimal.NewFromFloat(0.87894954884928), decimal.NewFromFloat(1.01333200284909), decimal.NewFromFloat(1.03490792577263), decimal.NewFromFloat(1.03490792577263),
				decimal.NewFromFloat(1.03490792577263), decimal.NewFromFloat(1.03490792577263), decimal.NewFromFloat(1.03490792577263)},
			Country:    "TR",
			CampaignID: "u6jyfnon-1f1d-4d86-91e9-1e00167939cb",
			UsersCount: 1,
		},
	}

	revenues, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error parsing CSV: %v", err)
	}
	assert.Equal(t, expectedRevenues, revenues)
}

func TestCSVParser_Parse_InvalidFile(t *testing.T) {
	parser := CSVParser{
		Path: "invalid_file.csv",
	}

	revenues, err := parser.Parse()
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, revenues)
}

func TestCSVParser_Parse_EmptyCSV(t *testing.T) {
	// Test parsing an empty CSV file
	emptyCSVData := ""

	// Create a temporary empty CSV file
	tempEmptyCSVFilePath, err := createTempCSVFile(emptyCSVData)
	if err != nil {
		t.Fatalf("Failed to create temp empty CSV file: %v", err)
	}
	defer removeTempCSVFile(tempEmptyCSVFilePath)

	parser := CSVParser{
		Path: tempEmptyCSVFilePath,
	}

	revenues, err := parser.Parse()
	assert.Error(t, err)
	assert.Equal(t, "EOF", err.Error())
	assert.Nil(t, revenues)
}

func TestCSVParser_Parse_InvalidCSV_WrongSeparator(t *testing.T) {
	// Test parsing an empty CSV file with only header
	emptyCSVData := `UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7
		1;81855ad8-681d-4d86-91e9-1e00167939cb;TR;1.54996978744822;2.2526636056983;2.29863633234526;2.88400864327196;3.6960018085883;5.7144365110237;0`

	// Create a temporary empty CSV file
	tempEmptyCSVFilePath, err := createTempCSVFile(emptyCSVData)
	if err != nil {
		t.Fatalf("Failed to create temp empty CSV file: %v", err)
	}
	defer removeTempCSVFile(tempEmptyCSVFilePath)

	parser := CSVParser{
		Path: tempEmptyCSVFilePath,
	}

	revenues, err := parser.Parse()
	assert.Error(t, err)
	assert.Nil(t, revenues)
}

func TestConvertCSVRecordToRevenues_InvalidLTV(t *testing.T) {
	// Test conversion with invalid LTV values in CSV record
	invalidLTVRecord := []string{"1", "campaign_id", "US", "1.5499697874482206", "2.252663605698363", "2.2986363323452683",
		"2.8840086432719603", "3.696001808588305", "invalid_ltv", "4.696001808588305"}

	revenues, err := convertCSVRecordToRevenues(invalidLTVRecord)
	assert.Error(t, err)
	assert.Equal(t, "can't convert invalid_ltv to decimal", err.Error())
	assert.Nil(t, revenues)
}

func TestConvertCSVRecordToRevenues_MissingFields(t *testing.T) {
	// Test conversion with missing fields in CSV record
	missingFieldsRecord := []string{"1", "campaign_id"}

	revenues, err := convertCSVRecordToRevenues(missingFieldsRecord)
	assert.Error(t, err, ErrNotEnoughFields)
	assert.Nil(t, revenues)
}
