package config

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"ohlc/util"
)

func RedpandaClient() *kgo.Client {
	bootstrapServer := util.GetEnvString("REDPANDA_SERVER", "127.0.0.1:9092")
	topic := util.GetEnvString("TOPIC", "ohlc")

	seeds := []string{bootstrapServer}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("ohlc-engine"),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		panic(err)
	}
	//defer cl.Close()

	return cl
}
