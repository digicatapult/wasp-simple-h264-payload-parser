package kafka_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Shopify/sarama"

	"github.com/digicatapult/wasp-simple-h264-payload-parser/kafka"
	"github.com/digicatapult/wasp-simple-h264-payload-parser/testdata"
)

func TestProducerConsumer(t *testing.T) {
	topic := "videos"
	kafkaBrokers := []string{"localhost:9092"}

	receivedMsgs := make(chan []byte)

	// setup consumer
	go kafka.SetupConsumer(kafkaBrokers, topic, receivedMsgs)

	// setup producer
	var (
		producer sarama.SyncProducer
		err      error
	)

	if producer, err = sarama.NewSyncProducer(kafkaBrokers, nil); err != nil {
		t.Fatalf("problem initiating producer: %s", err)
	}

	// send 5 messages?
	msg := testdata.KafkaMessage{
		Ingest:    "",
		IngestID:  "",
		Timestamp: "",
		Payload:   "",
		Metadata:  "",
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	testdata.SendMessage(producer, topic, "hello", msg)

	go func() {
		for {
			select {
			case msg := <-receivedMsgs:
				fmt.Println(msg)
				wg.Done()
			}
		}
	}()

	wg.Wait()
}
