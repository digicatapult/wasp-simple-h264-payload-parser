package main

import (
	"github.com/digicatapult/wasp-simple-h264-payload-parser/kafka"
)

func main() {
	topicName := "videos"
	kafkaBrokers := []string{"localhost:9092"}

	receivedMsgs := make(chan []byte)

	kafka.SetupConsumer(kafkaBrokers, topicName, receivedMsgs)

	/// do something with the messages

	// 	// Trap SIGINT to trigger a shutdown.
	// 	signals := make(chan os.Signal, 1)
	// 	signal.Notify(signals, os.Interrupt)

	// 	consumed := 0
	// ReceivedLoop:
	// 	for {
	// 		select {
	// 		case msg := <-receivedMsgs:
	// 			fmt.Println(msg)
	// 		case <-signals:
	// 			break ReceivedLoop
	// 		}
	// 	}
}
