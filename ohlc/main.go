package main

import (
	"context"
	"fmt"
	"github.com/alphadose/haxmap"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"ohlc/config"
	"ohlc/eventprocessor"
	"ohlc/models"
)

func consumeEvent(ctx context.Context, c *kgo.Client, r *redis.Client, hmap *haxmap.Map[string, *models.OHLC]) {
	for {
		fetches := c.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			eventprocessor.EventProcessor(ctx, r, string(record.Value), hmap)
			c.CommitRecords(ctx, record)
		}
	}
}

func main() {
	log.Println("Starting OHLC Engine")
	ctx := context.Background()
	redpandaClient := config.RedpandaClient()
	redisClient := config.GetRedisClient()
	hmap := haxmap.New[string, *models.OHLC]()
	consumeEvent(ctx, redpandaClient, redisClient, hmap)
}
