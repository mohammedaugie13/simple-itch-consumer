package models

import (
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
	expected := "{\"type\":\"P\",\"price\":\"2\",\"quantity\":\"1\",\"stock_code\":\"BBCA\"}"
	if expected != string(evntPJson) {
		t.Errorf("Expected value %v", expected)
	}
}

func TestFromJson(t *testing.T) {
	var err error
	eventEString := "{\"type\":\"E\",\"executed_price\":\"2\",\"executed_quantity\":\"1\",\"stock_code\":\"BBCA\"}"
	evntE := &EventMessage{}
	err = evntE.FromJSON([]byte(eventEString))
	t.Logf("error %v", err)

	eventPString := "{\"type\":\"P\",\"price\":\"2\",\"quantity\":\"1\",\"stock_code\":\"BBCA\"}"
	evntP := &EventMessage{}
	err = evntP.FromJSON([]byte(eventPString))
	t.Logf("error %v", err)
	if evntE.Type != "E" {
		t.Errorf("Expected %v", "E")
	}

	if evntP.Type != "P" {
		t.Errorf("Expected %v", "P")
	}
}
