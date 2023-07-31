package eventprocessor

import (
	"ohlc/models"
	"ohlc/util"
	"testing"

	"github.com/alphadose/haxmap"
)

func DecimalBig(val string) *util.StandardBigDecimal {
	a, _ := util.NewDecimalFromString(val)
	return a
}

func TestNDJsonAdd(t *testing.T) {
	data := `{"type":"A","order_book":"102","price":"6550","stock_code":"ASII"}
	{"type":"A","order_book":"722","price":"8875","stock_code":"BBCA"}
	{"type":"A","order_book":"911","price":"4650","stock_code":"BBRI"}
	{"type":"A","order_book":"192","price":"22600","stock_code":"GGRM"}
	{"type":"A","order_book":"207","price":"930","stock_code":"HMSP"}
	{"type":"A","order_book":"1578","price":"9750","stock_code":"ICBP"}
	{"type":"A","order_book":"552","price":"4190","stock_code":"TLKM"}
	{"type":"A","order_book":"35","price":"4540","stock_code":"UNVR"}`

	events := ProcessEvent(data)
	for _, event := range events {
		if event.Type != "A" {
			t.Errorf("Expected %v", "A")
		}
	}
}

func TestNDJsonOHLC(t *testing.T) {
	data := `{"type":"P","executed_quantity":"5","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
{"type":"P","executed_quantity":"7","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
{"type":"P","executed_quantity":"1","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
{"type":"P","executed_quantity":"1","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
{"type":"P","executed_quantity":"33","order_book":"35","execution_price":"4530","stock_code":"UNVR"}`
	hmap := haxmap.New[string, *models.OHLC]()
	events := ProcessEvent(data)
	details, _ := CalculateOHLC(events, hmap)
	results, _ := details.Get("UNVR")
	expected := models.OHLC{
		Volume:        DecimalBig("47"),
		HighestPrice:  DecimalBig("4530"),
		ClosePrice:    DecimalBig("4530"),
		Value:         DecimalBig("212910"),
		OpenPrice:     DecimalBig("4530"),
		PreviousPrice: DecimalBig("0"),
		AveragePrice:  DecimalBig("4530"),
		LowestPrice:   DecimalBig("4530"),
	}

	eJSON, _ := results.ToJSON()
	expectedJSON, _ := expected.ToJSON()
	if string(eJSON) != string(expectedJSON) {
		t.Errorf("Error %v", string(expectedJSON))
	}
}

func TestNDJsonOHLC2(t *testing.T) {
	data := `{"type":"A","quantity":"0","price":"8000","stock_code":"BBCA"}
{"type":"P","quantity":"100","price":"8050","stock_code":"BBCA"}
{"type":"P","quantity":"500","price":"7950","stock_code":"BBCA"}
{"type":"A","quantity":"200","price":"8150","stock_code":"BBCA"}
{"type":"E","quantity":"300","price":"8100","stock_code":"BBCA"}
{"type":"A","quantity":"100","price":"8200","stock_code":"BBCA"}
`
	hmap := haxmap.New[string, *models.OHLC]()
	events := ProcessEvent(data)
	details, _ := CalculateOHLC(events, hmap)

	results, _ := details.Get("BBCA")
	expected := models.OHLC{
		Volume:        DecimalBig("900"),
		HighestPrice:  DecimalBig("8100"),
		ClosePrice:    DecimalBig("8100"),
		Value:         DecimalBig("7210000"),
		OpenPrice:     DecimalBig("8050"),
		PreviousPrice: DecimalBig("8000"),
		AveragePrice:  DecimalBig("8011"),
		LowestPrice:   DecimalBig("7950"),
	}

	eJSON, _ := results.ToJSON()
	expectedJSON, _ := expected.ToJSON()
	if string(eJSON) != string(expectedJSON) {
		t.Errorf("Error %v", string(expectedJSON))
	}
}

func TestNDJsonOHLC3(t *testing.T) {
	data := `{"type":"A","quantity":"0","price":"8000","stock_code":"BBCA"}
{"type":"P","quantity":"100","price":"8050","stock_code":"BBCA"}
{"type":"P","quantity":"500","price":"7950","stock_code":"BBCA"}
{"type":"A","quantity":"200","price":"8150","stock_code":"BBCA"}
{"type":"E","quantity":"300","price":"8100","stock_code":"BBCA"}
{"type":"A","quantity":"100","price":"8200","stock_code":"BBCA"}
`
	hmap := haxmap.New[string, *models.OHLC]()
	events := ProcessEvent(data)
	CalculateOHLC(events, hmap)

	// results, _ := details.Get("BBCA")

	data2 := `{"type":"P","quantity":"100","price":"8050","stock_code":"BBCA"}`
	events2 := ProcessEvent(data2)
	details2, _ := CalculateOHLC(events2, hmap)

	results2, _ := details2.Get("BBCA")
	expected := models.OHLC{
		Volume:        DecimalBig("100"),
		HighestPrice:  DecimalBig("8050"),
		ClosePrice:    DecimalBig("8050"),
		Value:         DecimalBig("805000"),
		OpenPrice:     DecimalBig("8050"),
		PreviousPrice: DecimalBig("8000"),
		AveragePrice:  DecimalBig("8050"),
		LowestPrice:   DecimalBig("8050"),
	}

	eJSON, _ := results2.ToJSON()
	expectedJSON, _ := expected.ToJSON()
	if string(eJSON) != string(expectedJSON) {
		t.Errorf("Error %v", string(expectedJSON))
	}
}
