package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/luispinto23/movieexample/rating/pkg/model"
)

func main() {
	fmt.Printf("Creating a Kafka producer")

	// Create new kafka producer with localhost as the address
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratingsData.json"
	fmt.Println("Reading rating events from file " + fileName)

	RatingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceRatingEvents(topic, producer, RatingEvents); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second
	fmt.Println("Waiting " + timeout.String() + "until all events get produced")

	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func produceRatingEvents(topic string, producer *kafka.Producer, events []model.RatingEvent) error {
	for _, RatingEvent := range events {
		encodedEvent, err := json.Marshal(RatingEvent)
		if err != nil {
			return err
		}

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(encodedEvent),
		}, nil)

		if err != nil {
			return err
		}
	}
	return nil
}
