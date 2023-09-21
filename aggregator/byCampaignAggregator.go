package aggregator

import "github.com/pklimuk/ltv-predictor/fileParser"

type ByCampaignAggregator struct{}

func (a ByCampaignAggregator) AggregateRevenues(revenues []fileParser.Revenues) (AggregatedRevenuesByKey, error) {
	var result AggregatedRevenuesByKey = make(map[string]AggregatedRevenues)
	for i := 0; i < len(revenues); i++ {
		rec := revenues[i]
		if ar, ok := result[rec.CampaignID]; !ok {
			result[rec.CampaignID] = AggregatedRevenues{
				Revenues:   rec.Revenues,
				UsersCount: rec.UsersCount,
			}
		} else {
			err := ar.addRevenues(rec.Revenues)
			if err != nil {
				return nil, err
			}
			ar.UsersCount += rec.UsersCount
			result[rec.CampaignID] = ar
		}
	}
	return result, nil
}

func (a ByCampaignAggregator) ConvertAggregatedByKeyRevenuesToLTVs(ar AggregatedRevenuesByKey) (AggregatedLTVsByKey, error) {
	return convertAggregatedByKeyRevenuesToLTVs(ar)
}
