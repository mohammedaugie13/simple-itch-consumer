package models

import (
	"encoding/json"
	"errors"
	"github.com/alphadose/haxmap"
	"log"
	"ohlc/util"
)

const (
	ADD      string = "A"
	EXECUTED string = "E"
	TRADE    string = "P"
)

type EventMessage struct {
	Type      string                   `json:"type"`
	Price     *util.StandardBigDecimal `json:"price"`
	Quantity  *util.StandardBigDecimal `json:"quantity"`
	StockCode string                   `json:"stock_code"`
}

func (em *EventMessage) UnmarshalJSON(data []byte) error {
	obj := struct {
		Type             string `json:"type"`
		Price            string `json:"price,omitempty"`
		ExecutionPrice   string `json:"execution_price,omitempty"`
		ExecutedQuantity string `json:"executed_quantity,omitempty"`
		Quantity         string `json:"quantity,omitempty"`
		StockCode        string `json:"stock_code"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		log.Printf("Error Load JSON")
		return err
	}
	var err error
	em.Type = obj.Type

	if obj.Price == "" {
		em.Price, err = util.NewDecimalFromString(obj.ExecutionPrice)
		if err != nil {
			return errors.New("Invalid Price")
		}
	} else {
		em.Price, err = util.NewDecimalFromString(obj.Price)
		if err != nil {
			return errors.New("Invalid Price")
		}
	}

	if obj.Quantity == "" {
		if obj.ExecutedQuantity != "" {
			em.Quantity, err = util.NewDecimalFromString(obj.ExecutedQuantity)
			if err != nil {
				return errors.New("Invalid Executed Quantity")
			}
		} else {
			em.Quantity = util.DecimalZero()

		}

	} else {
		em.Quantity, err = util.NewDecimalFromString(obj.Quantity)
		if err != nil {
			return errors.New("Invalid Quantity")
		}
	}

	em.StockCode = obj.StockCode

	return nil
}

func (em *EventMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Type             string `json:"type"`
			Price            string `json:"price,omitempty"`
			ExecutionPrice   string `json:"execution_price,omitempty"`
			ExecutedQuantity string `json:"executed_quantity,omitempty"`
			Quantity         string `json:"quantity,omitempty"`
			StockCode        string `json:"stock_code"`
		}{
			Type:      em.Type,
			Price:     em.Price.String(),
			Quantity:  em.Quantity.String(),
			StockCode: em.StockCode,
		})
}

func (em *EventMessage) ToJSON() ([]byte, error) {
	str, err := json.Marshal(em)
	return str, err
}

func (em *EventMessage) FromJSON(msg []byte) error {
	err := json.Unmarshal(msg, em)
	if err != nil {
		return err
	}
	return nil
}

type OHLC struct {
	PreviousPrice *util.StandardBigDecimal `json:"previous_price"`
	OpenPrice     *util.StandardBigDecimal `json:"open_price"`
	HighestPrice  *util.StandardBigDecimal `json:"highest_price"`
	LowestPrice   *util.StandardBigDecimal `json:"lowest_price"`
	ClosePrice    *util.StandardBigDecimal `json:"close_price"`
	AveragePrice  *util.StandardBigDecimal `json:"average_price"`
	Volume        *util.StandardBigDecimal `json:"volume"`
	Value         *util.StandardBigDecimal `json:"value"`
}

func (o *OHLC) UnmarshalJSON(data []byte) error {
	obj := struct {
		PreviousPrice string `json:"previous_price"`
		OpenPrice     string `json:"open_price"`
		HighestPrice  string `json:"highest_price"`
		LowestPrice   string `json:"lowest_price"`
		ClosePrice    string `json:"close_price"`
		AveragePrice  string `json:"average_price"`
		Volume        string `json:"volume"`
		Value         string `json:"value"`
	}{}
	var err error
	if err := json.Unmarshal(data, &obj); err != nil {
		log.Printf("Error Load JSON")
		return err
	}

	o.PreviousPrice, err = util.NewDecimalFromString(obj.PreviousPrice)
	if err != nil {
		return errors.New("Invalid Previous Price")
	}

	o.OpenPrice, err = util.NewDecimalFromString(obj.OpenPrice)
	if err != nil {
		return errors.New("Invalid Open Price")
	}

	o.HighestPrice, err = util.NewDecimalFromString(obj.HighestPrice)
	if err != nil {
		return errors.New("Invalid Highest Price")
	}
	o.LowestPrice, err = util.NewDecimalFromString(obj.LowestPrice)
	if err != nil {
		return errors.New("Invalid Lowest Price")
	}
	o.ClosePrice, err = util.NewDecimalFromString(obj.ClosePrice)
	if err != nil {
		return errors.New("Invalid Close Price")
	}

	o.AveragePrice, err = util.NewDecimalFromString(obj.AveragePrice)
	if err != nil {
		return errors.New("Invalid Average Price")
	}

	o.Volume, err = util.NewDecimalFromString(obj.Volume)
	if err != nil {
		return errors.New("Invalid Volume Price")
	}

	o.Value, err = util.NewDecimalFromString(obj.Value)
	if err != nil {
		return errors.New("Invalid Value Price")
	}

	return nil
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
			OpenPrice:     o.OpenPrice.String(),
			HighestPrice:  o.HighestPrice.String(),
			LowestPrice:   o.LowestPrice.String(),
			ClosePrice:    o.ClosePrice.String(),
			AveragePrice:  o.AveragePrice.String(),
			Volume:        o.Volume.String(),
			Value:         o.Value.String(),
			PreviousPrice: o.PreviousPrice.String(),
		})
}

func (o *OHLC) ToJSON() ([]byte, error) {
	str, err := json.Marshal(o)
	return str, err
}

func (o *OHLC) FromJSON(msg []byte) error {
	err := json.Unmarshal(msg, o)
	if err != nil {
		return err
	}
	return nil
}

func NewOHLC() *OHLC {
	return &OHLC{
		OpenPrice:     util.DecimalZero(),
		HighestPrice:  util.DecimalZero(),
		LowestPrice:   util.DecimalZero(),
		ClosePrice:    util.DecimalZero(),
		AveragePrice:  util.DecimalZero(),
		Volume:        util.DecimalZero(),
		Value:         util.DecimalZero(),
		PreviousPrice: util.DecimalZero(),
	}
}

type OHLCMap = haxmap.Map[string, *OHLC]
