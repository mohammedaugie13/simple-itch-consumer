package eventprocessor

import (
	"github.com/alphadose/haxmap"
	"github.com/stretchr/testify/assert"
	"ohlc/models"
	"ohlc/util"
	"testing"
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
		assert.Equal(t, event.Type, "A")
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
	for _, event := range events {
		//t.Logf("DETAILS %v", event)
		assert.Equal(t, event.Type, "P")
	}
	results, _ := details.Get("UNVR")
	assert.Equal(t, results.Volume, DecimalBig("47"))
	assert.Equal(t, results.HighestPrice, DecimalBig("4530"))
	assert.Equal(t, results.ClosePrice, DecimalBig("4530"))
	assert.Equal(t, results.Value, DecimalBig("212910"))
	assert.Equal(t, results.OpenPrice, DecimalBig("4530"))
	assert.Equal(t, results.PreviousPrice, DecimalBig("0"))

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
	assert.Equal(t, results.PreviousPrice, DecimalBig("8000"))
	assert.Equal(t, results.AveragePrice, DecimalBig("8011"))
	assert.Equal(t, results.HighestPrice, DecimalBig("8100"))
	assert.Equal(t, results.LowestPrice, DecimalBig("7950"))
	assert.Equal(t, results.ClosePrice, DecimalBig("8100"))
	assert.Equal(t, results.Value, DecimalBig("7210000"))

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
	details, _ := CalculateOHLC(events, hmap)

	results, _ := details.Get("BBCA")
	assert.Equal(t, results.PreviousPrice, DecimalBig("8000"))
	assert.Equal(t, results.AveragePrice, DecimalBig("8011"))
	assert.Equal(t, results.HighestPrice, DecimalBig("8100"))
	assert.Equal(t, results.LowestPrice, DecimalBig("7950"))
	assert.Equal(t, results.ClosePrice, DecimalBig("8100"))
	assert.Equal(t, results.Value, DecimalBig("7210000"))

	data2 := `{"type":"P","quantity":"100","price":"8050","stock_code":"BBCA"}`
	events2 := ProcessEvent(data2)
	details2, _ := CalculateOHLC(events2, hmap)

	results2, _ := details2.Get("BBCA")
	assert.Equal(t, results2.PreviousPrice, DecimalBig("8000"))
	assert.Equal(t, results2.OpenPrice, DecimalBig("8050"))

}
