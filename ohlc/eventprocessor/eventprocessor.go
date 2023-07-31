package eventprocessor

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"ohlc/models"
	"ohlc/util"
	"strings"
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
	for _, datum := range em {

		details, ok := hmap.Get(datum.StockCode)
		if !ok {
			newOHLC := models.NewOHLC()
			if datum.Type == models.ADD {
				if datum.Quantity.Equal(util.DecimalZero().Decimal) {
					newOHLC.PreviousPrice = datum.Price
					hmap.Set(datum.StockCode, newOHLC)
				}
				continue
			}
			newOHLC.Volume = newOHLC.Volume.Add(datum.Quantity)
			newOHLC.ClosePrice = datum.Price
			newOHLC.HighestPrice = datum.Price
			newOHLC.LowestPrice = datum.Price
			newOHLC.Value = newOHLC.Value.Add(datum.Price.Mul(datum.Quantity))
			averagePrice := newOHLC.Value.Div(newOHLC.Volume).IntPart()
			newOHLC.AveragePrice = util.NewDecimalFromInt(averagePrice)
			if datum.Quantity != util.DecimalZero() && !slices.Contains(foundOpenPrice, datum.StockCode) {
				// first occurrence
				newOHLC.OpenPrice = datum.Price
				seen = append(seen, datum.StockCode)
				hmap.Set(datum.StockCode, newOHLC)
			} else {
				newOHLC.PreviousPrice = datum.Price
				//seen = append(seen, datum.StockCode)
				hmap.Set(datum.StockCode, newOHLC)

			}
		} else {

			if datum.Type == models.ADD {
				if datum.Quantity == util.DecimalZero() {
					details.PreviousPrice = datum.Price
					hmap.Set(datum.StockCode, details)

				}
				continue
			}

			if slices.Contains(seen, datum.StockCode) {
				if details.LowestPrice.GreaterThan(datum.Price.Decimal) {
					details.LowestPrice = datum.Price
				}
				if details.HighestPrice.LessThan(datum.Price.Decimal) {
					details.HighestPrice = datum.Price
				}

				details.Volume = details.Volume.Add(datum.Quantity)
				details.Value = details.Value.Add(datum.Price.Mul(datum.Quantity))
				averagePrice := details.Value.Div(details.Volume).IntPart()
				details.AveragePrice = util.NewDecimalFromInt(averagePrice)
			} else {
				details.Volume = datum.Quantity
				details.HighestPrice = datum.Price
				details.LowestPrice = datum.Price
				details.Value = datum.Price.Mul(datum.Quantity)
				averagePrice := details.Value.Div(details.Volume).IntPart()
				details.AveragePrice = util.NewDecimalFromInt(averagePrice)
			}
			details.ClosePrice = datum.Price

			if !slices.Contains(foundOpenPrice, datum.StockCode) {
				details.OpenPrice = datum.Price
				seen = append(foundOpenPrice, datum.StockCode)
			}
			hmap.Set(datum.StockCode, details)

		}

	}
	return hmap, seen
}

func SaveToRedis(ctx context.Context, r *redis.Client, seen []string, hmap *models.OHLCMap) {
	for _, datum := range seen {
		details, _ := hmap.Get(datum)
		detailsJson, _ := details.ToJSON()
		r.HSet(ctx, datum, detailsJson)
	}
}
