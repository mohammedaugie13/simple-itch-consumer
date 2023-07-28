package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToJson(t *testing.T) {
	eventE := EventMessage{Type: "E", ExecutedQuantity: "1", ExecutedPrice: "2", StockCode: "BBCA"}
	evntEJson, _ := eventE.ToJSON()

	assert.Equal(t, "{\"type\":\"E\",\"executed_price\":\"2\",\"executed_quantity\":\"1\",\"stock_code\":\"BBCA\"}", string(evntEJson))

	eventP := EventMessage{Type: "P", Quantity: "1", Price: "2", StockCode: "BBCA"}
	evntPJson, _ := eventP.ToJSON()
	assert.Equal(t, "{\"type\":\"P\",\"price\":\"2\",\"quantity\":\"1\",\"stock_code\":\"BBCA\"}", string(evntPJson))

}

func TestFromJson(t *testing.T) {
	eventEString := "{\"type\":\"E\",\"executed_price\":\"2\",\"executed_quantity\":\"1\",\"stock_code\":\"BBCA\"}"
	evntE := &EventMessage{}
	evntE.FromJSON([]byte(eventEString))

	eventPString := "{\"type\":\"P\",\"price\":\"2\",\"quantity\":\"1\",\"stock_code\":\"BBCA\"}"
	evntP := &EventMessage{}
	evntP.FromJSON([]byte(eventPString))

	assert.Equal(t, "E", evntE.Type)
	assert.Equal(t, "P", evntP.Type)
	assert.Equal(t, "2", evntE.ExecutedPrice)
	assert.Equal(t, "2", evntP.Quantity)

}
