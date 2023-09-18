package fileParser

import (
	"encoding/csv"
	"os"

	"github.com/shopspring/decimal"
)

type CSVParser struct {
	Path string
}

func (p CSVParser) Parse() ([]Revenues, error) {
	records, err := parseCSV(p.Path)
	if err != nil {
		return nil, err
	}
	revenues := make([]Revenues, 0, len(records))
	for _, record := range records {
		revenue, err := convertCSVRecordToRevenues(record)
		if err != nil {
			return nil, err
		}
		revenues = append(revenues, revenue)
	}
	return revenues, nil
}

func parseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// Read and skip header row
	_, err = csvReader.Read()
	if err != nil {
		return nil, err
	}
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func convertCSVRecordToRevenues(record []string) (Revenues, error) {
	campaignID := record[1]
	country := record[2]
	var ltv = make([]decimal.Decimal, 0, len(record)-3)
	for i := 3; i < len(record); i++ {
		ltvValue, err := decimal.NewFromString(record[i])
		if err != nil {
			return Revenues{}, err
		}
		ltv = append(ltv, ltvValue)
		normalizeLtv(&ltv)
	}
	return Revenues{Revenues: ltv, Country: country, CampaignID: campaignID, UsersCount: 1}, nil
}

// normalizeLtv replaces zero values with the previous non-zero value
func normalizeLtv(ltv *[]decimal.Decimal) {
	prevLtv := (*ltv)[0]
	for i := 1; i < len(*ltv); i++ {
		if (*ltv)[i].Equal(decimal.Zero) {
			(*ltv)[i] = prevLtv
		}
		prevLtv = (*ltv)[i]
	}
}
