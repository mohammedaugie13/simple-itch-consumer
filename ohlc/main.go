package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"ohlc/config"
)

func consumeEvent(ctx context.Context, c *kgo.Client) {
	for {
		fetches := c.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			//log.Printf("Value %v", string(record.Value))
			c.CommitRecords(ctx, record)
		}
	}
}

func main() {
	log.Println("Starting OHLC Engine")
	ctx := context.Background()
	redpandaClient := config.RedpandaClient()
	consumeEvent(ctx, redpandaClient)
}
