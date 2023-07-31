package models

import "encoding/json"

type OHLC struct {
	PreviousPrice string `json:"previous_price"`
	OpenPrice     string `json:"open_price"`
	HighestPrice  string `json:"highest_price"`
	LowestPrice   string `json:"lowest_price"`
	ClosePrice    string `json:"close_price"`
	AveragePrice  string `json:"average_price"`
	Volume        string `json:"volume"`
	Value         string `json:"value"`
}

func (o *OHLC) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			PreviousPrice string `json:"previous_price"`
			OpenPrice     string `json:"open_price"`
			HighestPrice  string `json:"highest_price"`
			LowestPrice   string `json:"lowest_price"`
			ClosePrice    string `json:"close_price"`
			AveragePrice  string `json:"average_price"`
			Volume        string `json:"volume"`
			Value         string `json:"value"`
		}{
			OpenPrice:     o.OpenPrice,
			HighestPrice:  o.HighestPrice,
			LowestPrice:   o.LowestPrice,
			ClosePrice:    o.ClosePrice,
			AveragePrice:  o.AveragePrice,
			Volume:        o.Volume,
			Value:         o.Value,
			PreviousPrice: o.PreviousPrice,
		})
}

func (o *OHLC) ToJSON() ([]byte, error) {
	str, err := json.Marshal(o)
	return str, err
}
