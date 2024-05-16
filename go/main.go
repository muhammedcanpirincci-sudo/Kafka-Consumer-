package main

import (
    "context"
    "math/rand"
    "time"

    "github.com/cenkalti/backoff/v4"
    "github.com/segmentio/kafka-go"
    kafka_retry "yourmodule/name/kafka_retry"
)


func main() {
    rand.Seed(time.Now().UnixNano())
    ctx := context.Background()
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "products-topic",
        GroupID: "your-consumer-group",
        StartOffset: kafka.LastOffset,
        // Additional configuration may be required depending on your Kafka setup
    })

    options := kafka_retry.ConsumerWithRetryOptions{
        Handler: &kafka_retry.SimpleHandler{},
        Reader:  reader,
        MaxRetries: 3,
        ProcessRetryQueue: make(chan kafka_retry.MessageWithIndex), // Corrected type
        ReadRetryQueue:  make(chan kafka_retry.MessageWithIndex), // Corrected type
        Backoff:    backoff.NewExponentialBackOff(),
    }
    kafka_retry.NewConsumerWithRetry(ctx, &options)
}
