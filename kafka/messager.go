package kafka

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

// Messager implements kafka message functionality
type Messager struct {
	sp sarama.SyncProducer
	c  sarama.Consumer
}

// NewMessager will instantiate an instance using the producer provided
func NewMessager() *Messager {
	return &Messager{}
}

// WithConsumer will set the consumer instance on the messager
func (m *Messager) WithConsumer(c sarama.Consumer) *Messager {
	m.c = c
	return m
}

// WithProducer will set the producer instance on the messager
func (m *Messager) WithProducer(p sarama.SyncProducer) *Messager {
	m.sp = p
	return m
}

// SendMessage can send a message to the
func (m *Messager) SendMessage(topic, mKey string, mValue []byte) error {
	if m.sp == nil {
		panic("nil producer")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(mKey),
		Value: sarama.StringEncoder(mValue),
	}

	partition, offset, err := m.sp.SendMessage(msg)
	if err != nil {
		zap.S().Errorf("error sending msg %s - %s (%d, %d", msg.Key, err, partition, offset)

		return err
	}

	zap.S().Debugf("Message sent to partition %d, offset %d", partition, offset)
	return nil
}

// Listen will listen indefinitely to the kafka bus for messages on a specific topic
func (m *Messager) Listen(topic string, received chan []byte) {
	if m.c == nil {
		panic("nil consumer")
	}
	partitionConsumer, err := m.c.ConsumePartition(topic, 0, 0)
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
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
			received <- msg.Value
		case <-signals:
			return
		}
	}
}
