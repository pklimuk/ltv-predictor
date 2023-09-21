package fileParser

import (
	"encoding/json"
	"io"
	"os"

	"github.com/shopspring/decimal"
)

type JSONParser struct {
	Path string
}

type jsonData struct {
	CampaignID string          `json:"CampaignId"`
	Country    string          `json:"Country"`
	Ltv1       decimal.Decimal `json:"Ltv1"`
	Ltv2       decimal.Decimal `json:"Ltv2"`
	Ltv3       decimal.Decimal `json:"Ltv3"`
	Ltv4       decimal.Decimal `json:"Ltv4"`
	Ltv5       decimal.Decimal `json:"Ltv5"`
	Ltv6       decimal.Decimal `json:"Ltv6"`
	Ltv7       decimal.Decimal `json:"Ltv7"`
	Users      int64           `json:"Users"`
}

func (p JSONParser) Parse() ([]Revenues, error) {
	data, err := parseJSONFile(p.Path)
	if err != nil {
		return nil, err
	}
	revenues := make([]Revenues, 0, len(data))
	for _, rec := range data {
		revenues = append(revenues, convertJSONDataToRevenue(rec))
	}
	return revenues, nil
}

func convertJSONDataToRevenue(d jsonData) Revenues {
	ltvs := []decimal.Decimal{d.Ltv1, d.Ltv2, d.Ltv3, d.Ltv4, d.Ltv5, d.Ltv6, d.Ltv7}
	for i := 0; i < len(ltvs); i++ {
		ltvs[i] = ltvs[i].Mul(decimal.NewFromInt(d.Users))
	}
	return Revenues{Revenues: ltvs, Country: d.Country, CampaignID: d.CampaignID, UsersCount: d.Users}
}

func parseJSONFile(path string) ([]jsonData, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var data []jsonData
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
