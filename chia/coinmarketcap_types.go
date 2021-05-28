package chia

import (
	"time"
)

type ChiaCryptocurrencyMapResponse struct {
	Data []struct {
		ID                  int         `json:"id"`
		Rank                int         `json:"rank"`
		Name                string      `json:"name"`
		Symbol              string      `json:"symbol"`
		Slug                string      `json:"slug"`
		IsActive            int         `json:"is_active"`
		FirstHistoricalData time.Time   `json:"first_historical_data"`
		LastHistoricalData  time.Time   `json:"last_historical_data"`
		Platform            interface{} `json:"platform"`
	} `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}
type ChiaTableDbResponse struct {
	UpdateId         int64
	ChiaPrice        float64
	PercentChange1H  float64
	PercentChange24h float64
	TotalSupply      float64
}
type CoinMarketCapSymbolResponse struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data struct {
		Num9258 struct {
			ID                int           `json:"id"`
			Name              string        `json:"name"`
			Symbol            string        `json:"symbol"`
			Slug              string        `json:"slug"`
			NumMarketPairs    int           `json:"num_market_pairs"`
			DateAdded         time.Time     `json:"date_added"`
			Tags              []interface{} `json:"tags"`
			MaxSupply         interface{}   `json:"max_supply"`
			CirculatingSupply int           `json:"circulating_supply"`
			TotalSupply       int           `json:"total_supply"`
			IsActive          int           `json:"is_active"`
			Platform          interface{}   `json:"platform"`
			CmcRank           int           `json:"cmc_rank"`
			IsFiat            int           `json:"is_fiat"`
			LastUpdated       time.Time     `json:"last_updated"`
			Quote             struct {
				Usd struct {
					Price            float64   `json:"price"`
					Volume24H        float64   `json:"volume_24h"`
					PercentChange1H  float64   `json:"percent_change_1h"`
					PercentChange24H float64   `json:"percent_change_24h"`
					PercentChange7D  float64   `json:"percent_change_7d"`
					PercentChange30D int       `json:"percent_change_30d"`
					PercentChange60D int       `json:"percent_change_60d"`
					PercentChange90D int       `json:"percent_change_90d"`
					MarketCap        int       `json:"market_cap"`
					LastUpdated      time.Time `json:"last_updated"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"9258"`
	} `json:"data"`
}
