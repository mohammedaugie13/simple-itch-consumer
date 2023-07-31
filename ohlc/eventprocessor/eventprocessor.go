package eventprocessor

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"ohlc/models"
	"ohlc/util"
	"strings"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
)

func EventProcessor(ctx context.Context, r *redis.Client, event string, hmap *models.OHLCMap) {
	events := ProcessEvent(event)
	_, seen := CalculateOHLC(events, hmap)
	SaveToRedis(ctx, r, seen, hmap)
}

func ProcessEvent(event string) []models.EventMessage {
	d := json.NewDecoder(strings.NewReader(event))
	var events []models.EventMessage
	for {
		var v models.EventMessage
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		events = append(events, v)
	}
	return events
}

func CalculateOHLC(em []models.EventMessage, hmap *models.OHLCMap) (*models.OHLCMap, []string) {
	var seen []string
	var foundOpenPrice []string
	zeroState := models.NewOHLC()
	for _, datum := range em {
		details, _ := hmap.GetOrSet(datum.StockCode, zeroState)
		if datum.Type == models.ADD {
			if datum.Quantity.Equal(util.DecimalZero().Decimal) {
				details.PreviousPrice = datum.Price
				hmap.Set(datum.StockCode, details)
			}
			continue
		}
		if slices.Contains(seen, datum.StockCode) {
			SetLowestAndHighestPrice(datum, details)
			details.Volume = details.Volume.Add(datum.Quantity)
			details.Value = details.Value.Add(datum.Price.Mul(datum.Quantity))
			details.AveragePrice = util.NewDecimalFromInt(details.Value.Div(details.Volume).IntPart())
		} else {
			details.Volume = datum.Quantity
			details.HighestPrice = datum.Price
			details.LowestPrice = datum.Price
			details.Value = datum.Price.Mul(datum.Quantity)
			details.AveragePrice = util.NewDecimalFromInt(details.Value.Div(details.Volume).IntPart())
			seen = append(seen, datum.StockCode)
		}
		details.ClosePrice = datum.Price
		if !slices.Contains(foundOpenPrice, datum.StockCode) {
			UpdateOpenPrice(&foundOpenPrice, datum, details)
		}
		hmap.Set(datum.StockCode, details)
	}
	return hmap, seen
}

func UpdateOpenPrice(foundOpenPrices *[]string, event models.EventMessage, ohlc *models.OHLC) {
	ohlc.OpenPrice = event.Price
	*foundOpenPrices = append(*foundOpenPrices, event.StockCode)
}

func SetLowestAndHighestPrice(event models.EventMessage, ohlc *models.OHLC) {
	if ohlc.LowestPrice.GreaterThan(event.Price.Decimal) {
		ohlc.LowestPrice = event.Price
	}
	if ohlc.HighestPrice.LessThan(event.Price.Decimal) {
		ohlc.HighestPrice = event.Price
	}
}

func SaveToRedis(ctx context.Context, r *redis.Client, seen []string, hmap *models.OHLCMap) {
	log.Println("Save to Redis")
	log.Printf("Seen %v", seen)
	for _, datum := range seen {
		details, _ := hmap.Get(datum)
		detailsJSON, _ := details.ToJSON()
		cmd := r.Set(ctx, datum, string(detailsJSON), 0).Err()
		log.Printf("Success save to redis %v", cmd)
	}
}
