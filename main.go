package main

import (
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"github.com/digicatapult/wasp-simple-h264-payload-parser/kafka"
	"github.com/digicatapult/wasp-simple-h264-payload-parser/wasp"
)

func main() {
	inTopicName := "videos"
	outTopicName := "videos-transformed"
	kafkaBrokers := []string{"localhost:9092"}

	consumer, err := sarama.NewConsumer(kafkaBrokers, nil)
	if err != nil {
		panic(err)
	}

	producer, err := sarama.NewSyncProducer(kafkaBrokers, nil)
	if err != nil {
		zap.S().Fatalf("problem initiating producer: %s", err)
	}

	relayer := kafka.NewMessager().
		WithConsumer(consumer).
		WithProducer(producer).
		WithRelayTransform(wasp.TransformVideoMessages)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	relayer.Relay(inTopicName, outTopicName, stop)
}
