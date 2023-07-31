package main

import (
	"encoding/json"
	"fmt"
)

type ChangeRecord struct {
	StockCode string
	Price     int64
	Quantity  int64
}

type IndexMember struct {
	StockCode string
	IndexCode string
}

type Summary struct {
	StockCode string   `json:"stock_code"`
	IndexCode []string `json:"index_code"`
	Open      int64    `json:"open"`
	High      int64    `json:"high"`
	Low       int64    `json:"low"`
	Close     int64    `json:"close"`
	Prev      int64    `json:"prev"`
}

func ohlc(records map[string]*Summary, incomingData []ChangeRecord, indexes []IndexMember) {
	for _, datum := range incomingData {
		details, isFound := records[datum.StockCode]
		if !isFound {
			details = &Summary{
				StockCode: datum.StockCode,
				IndexCode: GetListedIndex(datum.StockCode, indexes),
			}
			records[datum.StockCode] = details
		}
		UpdateRecord(details, datum)
	}
}

func GetListedIndex(stockCode string, indexes []IndexMember) []string {
	var listedIndex []string
	for _, details := range indexes {
		if details.StockCode == stockCode {
			listedIndex = append(listedIndex, details.IndexCode)
		}
	}
	return listedIndex
}

func UpdateRecord(summary *Summary, incomingData ChangeRecord) {
	switch {
	case incomingData.Quantity == 0:
		summary.Prev = incomingData.Price
	case incomingData.Quantity > 0 && summary.Open == 0:
		summary.Open = incomingData.Price
	default:
		summary.Close = incomingData.Price
		if summary.High < incomingData.Price {
			summary.High = incomingData.Price
		}
		if summary.Low > incomingData.Price {
			summary.Low = incomingData.Price
		}
	}
}

var (
	incomingData = []ChangeRecord{
		{
			StockCode: "BBCA",
			Price:     8783,
			Quantity:  0,
		},
		{
			StockCode: "BBRI",
			Price:     3233,
			Quantity:  0,
		},
		{
			StockCode: "ASII",
			Price:     1223,
			Quantity:  0,
		},
		{
			StockCode: "GOTO",
			Price:     321,
			Quantity:  0,
		},

		{
			StockCode: "BBCA",
			Price:     8780,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3230,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1220,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     320,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8800,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3300,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1300,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     330,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8600,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3100,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1100,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     310,
			Quantity:  1,
		},

		{
			StockCode: "BBCA",
			Price:     8785,
			Quantity:  1,
		},
		{
			StockCode: "BBRI",
			Price:     3235,
			Quantity:  1,
		},
		{
			StockCode: "ASII",
			Price:     1225,
			Quantity:  1,
		},
		{
			StockCode: "GOTO",
			Price:     325,
			Quantity:  1,
		},
	}
	indexes = []IndexMember{
		{
			StockCode: "BBCA",
			IndexCode: "IHSG",
		},
		{
			StockCode: "BBRI",
			IndexCode: "IHSG",
		},
		{
			StockCode: "ASII",
			IndexCode: "IHSG",
		},
		{
			StockCode: "GOTO",
			IndexCode: "IHSG",
		},
		{
			StockCode: "BBCA",
			IndexCode: "LQ45",
		},
		{
			StockCode: "BBRI",
			IndexCode: "LQ45",
		},
		{
			StockCode: "ASII",
			IndexCode: "LQ45",
		},
		{
			StockCode: "GOTO",
			IndexCode: "LQ45",
		},
		{
			StockCode: "BBCA",
			IndexCode: "KOMPAS100",
		},
		{
			StockCode: "BBRI",
			IndexCode: "KOMPAS100",
		},
	}
)

func main() {
	summary := make(map[string]*Summary)
	ohlc(summary, incomingData, indexes)

	for _, v := range summary {
		jss, _ := json.Marshal(v)
		fmt.Println("summary : ", string(jss))
	}
}
