package kafkax

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/51xpx/pkgx/pkg/threading"
)

// Consumer kafka Consumer
type Consumer struct {
	Brokers []string
	Topic   string
	GroupID string
	reader  *kafka.Reader
}

// NewConsumer new kafka consumer
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
		reader:  r,
	}
}

// Close close client connection
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// Consume consume message from kafka
func (c *Consumer) Consume(fn func(message []byte)) {
	threading.GoSafe(func() {
		for {
			m, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			fn(m.Value)
		}
	})
}
