package fileParser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/shopspring/decimal"
)

// Constants for CSV file
const (
	fieldsNumber    = 10
	userIDIndex     = 0
	campaignIDIndex = 1
	countryIndex    = 2
	startLtvIndex   = 3
)

var (
	ErrCantReadHeader  = errors.New("can't read header row")
	ErrNotEnoughFields = errors.New("not enough fields in the record")
)

type CSVParser struct {
	Path string
}

func (p CSVParser) Parse() ([]Revenues, error) {
	records, err := parseCSV(p.Path)
	if err != nil {
		return nil, fmt.Errorf(ErrParsingError.Error(), err)
	}
	revenues := make([]Revenues, 0, len(records))
	for _, record := range records {
		revenue, err := convertCSVRecordToRevenues(record)
		// If there is an error, we just skip the record and log it, to not break the whole process
		if err != nil {
			log.Printf("Record %v contains errors(%v) and could not be processed.", record, err)
			continue
		}
		revenues = append(revenues, *revenue)
	}
	return revenues, nil
}

func parseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(ErrCantOpenFile.Error(), path)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// Read and skip header row
	_, err = csvReader.Read()
	if err != nil {
		return nil, ErrCantReadHeader
	}
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func convertCSVRecordToRevenues(record []string) (*Revenues, error) {
	if len(record) != fieldsNumber {
		return nil, ErrNotEnoughFields
	}
	campaignID := record[campaignIDIndex]
	country := record[countryIndex]
	var ltv = make([]decimal.Decimal, 0, len(record)-startLtvIndex)
	for i := 3; i < len(record); i++ {
		ltvValue, err := decimal.NewFromString(record[i])
		if err != nil {
			return nil, err
		}
		ltv = append(ltv, ltvValue)
		normalizeLtv(ltv)
	}
	return &Revenues{Revenues: ltv, Country: country, CampaignID: campaignID, UsersCount: 1}, nil
}

// normalizeLtv replaces zero values with the previous non-zero value
func normalizeLtv(ltv []decimal.Decimal) {
	prevLtv := (ltv)[0]
	for i := 1; i < len(ltv); i++ {
		if (ltv)[i].Equal(decimal.Zero) {
			(ltv)[i] = prevLtv
		}
		prevLtv = (ltv)[i]
	}
}
