package models

import (
	"encoding/json"
	"log"
)

const (
	A string = "A"
	E string = "E"
	P string = "P"
)

type EventMessage struct {
	Type             string `json:"type"`
	Price            string `json:"price,omitempty"`
	ExecutedPrice    string `json:"executed_price,omitempty"`
	ExecutedQuantity string `json:"executed_quantity,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	StockCode        string `json:"stock_code"`
}

func (em *EventMessage) UnmarshalJSON(data []byte) error {
	obj := struct {
		Type             string `json:"type"`
		Price            string `json:"price,omitempty"`
		ExecutedPrice    string `json:"executed_price,omitempty"`
		ExecutedQuantity string `json:"executed_quantity,omitempty"`
		Quantity         string `json:"quantity,omitempty"`
		StockCode        string `json:"stock_code"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		log.Printf("Error Load JSON")
		return err
	}
	em.Type = obj.Type
	em.Price = obj.Price
	em.ExecutedPrice = obj.ExecutedPrice
	em.Quantity = obj.ExecutedQuantity
	em.StockCode = obj.StockCode

	return nil
}

func (em *EventMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Type             string `json:"type"`
			Price            string `json:"price,omitempty"`
			ExecutedPrice    string `json:"executed_price,omitempty"`
			ExecutedQuantity string `json:"executed_quantity,omitempty"`
			Quantity         string `json:"quantity,omitempty"`
			StockCode        string `json:"stock_code"`
		}{
			Type:             em.Type,
			Price:            em.Price,
			ExecutedQuantity: em.ExecutedQuantity,
			Quantity:         em.Quantity,
			StockCode:        em.StockCode,
			ExecutedPrice:    em.ExecutedPrice,
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
