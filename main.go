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

func main() {
	inTopicName := util.GetEnv(util.InTopicNameKey, "payloads.simpleH264")
	outTopicName := util.GetEnv(util.OutTopicNameKey, "video")
	kafkaBrokersVar := util.GetEnv(util.KafkaBrokersKey, "localhost:9092")
	kafkaBrokers := strings.Split(kafkaBrokersVar, ",")

	cfg := zap.NewDevelopmentConfig()
	if util.GetEnv(util.EnvKey, "") == "production" {
		cfg = zap.NewProductionConfig()

		lvl, err := zap.ParseAtomicLevel(util.GetEnv(util.LogLevelKey, "debug"))
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
		syncErr := logger.Sync()
		if err != nil {
			log.Printf("error whilst syncing zap: %s\n", syncErr)
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

	err = relayer.Relay(inTopicName, outTopicName, stop)
	if err != nil {
		zap.S().Error(err)
	}

	select {
	case <-stop:
		zap.S().Info("closing down")
	}
}
