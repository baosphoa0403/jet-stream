package util

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNat() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	return nc
}

func CreateStream(js nats.JetStreamContext, config nats.StreamConfig) {
	jsInfo, err := js.AddStream(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stream EVENTS created with limits", jsInfo)
}

func CheckPendingAcks(js nats.JetStreamContext, streamName string, consumerName string) {
	consumerInfo, err := js.ConsumerInfo(streamName, consumerName)
	if err != nil {
		log.Printf("Error fetching consumer info: %v", err)
		return
	}

	// Display consumer info including unacknowledged messages
	fmt.Printf("Consumer Name: %s\n", consumerInfo.Name)
	fmt.Printf("Pending Messages to Acknowledge: %d\n", consumerInfo.NumAckPending)
	fmt.Printf("Total Messages Delivered: %d\n", consumerInfo.Delivered.Stream)
	fmt.Printf("Messages Redelivered: %d\n", consumerInfo.NumRedelivered)
}

func CheckInfoStreamExist(js nats.JetStreamContext, streamName string) bool {
	_, err := js.StreamInfo(streamName)
	return err == nil
}

func CreateStreamDeadLetterQueue(js nats.JetStreamContext) {
	streamName := "METRICS_DEAD_LETTER_QUEUE"

	if CheckInfoStreamExist(js, streamName) {
		log.Println("Stream already exists:", streamName)
	} else {
		config := nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{"dlq.>"},
			Retention: nats.WorkQueuePolicy, // Retain until acknowledged
			Storage:   nats.FileStorage,
		}

		CreateStream(js, config)

		fmt.Println("Create Stream Dead Letter Queue")
	}

}
