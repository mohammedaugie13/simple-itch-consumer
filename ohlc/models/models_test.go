package models

import (
	"github.com/stretchr/testify/assert"
	"ohlc/util"
	"testing"
)

func DecimalBig(val string) *util.StandardBigDecimal {
	a, _ := util.NewDecimalFromString(val)
	return a
}

func TestToJson(t *testing.T) {

	eventP := EventMessage{Type: "P", Quantity: DecimalBig("1"), Price: DecimalBig("2"), StockCode: "BBCA"}
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
	assert.Equal(t, DecimalBig("1"), evntE.Quantity)
	assert.Equal(t, DecimalBig("1"), evntP.Quantity)

}
