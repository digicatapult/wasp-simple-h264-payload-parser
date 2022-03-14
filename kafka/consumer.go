package kafka

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

// KafkaMessage defines the message structure
type KafkaMessage struct {
	Ingest    string `json:"ingest"`
	IngestID  string `json:"ingestId"`
	Timestamp string `json:"timestamp"`
	Payload   string `json:"payload"`
	Metadata  string `json:"metadata"`
}

type KafkaMessage2 struct {
	ThingID   string
	Type      string
	Timestamp string
	Value     string
	Metadata  string
}

// SetupConsumer will create a consumer and start listening for messages
func SetupConsumer(kafkaBrokers []string, topicName string, msgChannel chan []byte) {
	consumer, err := sarama.NewConsumer(kafkaBrokers, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topicName, 0, 0)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
			msgChannel <- msg.Value
		case <-signals:
			break ConsumerLoop
		}
	}
}
