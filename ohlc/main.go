package main

import (
	"context"
	"fmt"
	"github.com/alphadose/haxmap"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"
	"log"
	"net"
	"ohlc/config"
	"ohlc/engine"
	"ohlc/eventprocessor"
	"ohlc/models"
	s "ohlc/models/pb"
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

func InitGRPCServer(hmap *haxmap.Map[string, *models.OHLC]) {
	log.Println("Start GRPC Service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	cs := engine.Server{Map: hmap}

	s.RegisterOHLCServiceServer(gs, &cs)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func main() {
	ctx := context.Background()
	redpandaClient := config.RedpandaClient()
	redisClient := config.GetRedisClient()
	hmap := haxmap.New[string, *models.OHLC]()
	go InitGRPCServer(hmap)
	log.Println("Starting OHLC Engine")

	consumeEvent(ctx, redpandaClient, redisClient, hmap)
}
