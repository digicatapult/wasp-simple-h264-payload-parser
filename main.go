package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"github.com/digicatapult/wasp-simple-h264-payload-parser/kafka"
	"github.com/digicatapult/wasp-simple-h264-payload-parser/util"
	"github.com/digicatapult/wasp-simple-h264-payload-parser/wasp"
)

const (
	// InTopicNameKey defines the environment variable key for InTopicName
	InTopicNameKey = "IN_TOPIC_NAME_KEY"

	// OutTopicNameKey defines the environment variable key for OutTopicName
	OutTopicNameKey = "OUT_TOPIC_NAME_KEY"

	// KafkaBrokersKey defines the environment variable key for KafkaBrokers
	KafkaBrokersKey = "KAFKA_BROKERS"
)

func main() {
	inTopicName := util.GetEnv(InTopicNameKey, "payloads.simpleH264")
	outTopicName := util.GetEnv(OutTopicNameKey, "video")
	kafkaBrokersVar := util.GetEnv(KafkaBrokersKey, "localhost:9092")
	kafkaBrokers := strings.Split(kafkaBrokersVar, ",")

	cfg := zap.NewDevelopmentConfig()
	if os.Getenv("ENV") == "production" {
		cfg = zap.NewProductionConfig()

		lvl, err := zap.ParseAtomicLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			panic("invalid log level")
		}

		log.Printf("setting level: %s", lvl.String())

		cfg.Level = lvl
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("error initializing the logger")
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf("error whilst syncing zap: %s\n", err)
		}
	}()

	zap.ReplaceGlobals(logger)

	sarama.Logger = util.SaramaZapLogger{}

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
