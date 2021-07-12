package kafkax

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer kafka producer
type Producer struct {
	Brokers []string
	Topic   string
	writer  *kafka.Writer
}

// NewProducer .
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}

	return &Producer{
		Brokers: brokers,
		Topic:   topic,
		writer:  w,
	}
}

// Close close client connection
func (c *Producer) Close() error {
	return c.writer.Close()
}

// ProduceStringList produce string messages to kafka, returning the number of bytes written
func (c *Producer) ProduceStringList(msgs []string) error {
	var kmsgs []kafka.Message
	for _, msg := range msgs {
		kmsgs = append(kmsgs, kafka.Message{Value: []byte(msg)})
	}
	return c.writer.WriteMessages(context.Background(), kmsgs...)
}

// ProduceString produce string messages to kafka, returning the number of bytes written
func (c *Producer) ProduceString(msg string) error {
	return c.writer.WriteMessages(context.Background(), kafka.Message{Value: []byte(msg)})
}

// ProduceBytesList produce bytes messages to kafka, returning the number of bytes written
func (c *Producer) ProduceBytesList(msgs [][]byte) error {
	var kmsgs []kafka.Message
	for _, msg := range msgs {
		kmsgs = append(kmsgs, kafka.Message{Value: msg})
	}
	return c.writer.WriteMessages(context.Background(), kmsgs...)
}

// ProduceBytes produce bytes messages to kafka, returning the number of bytes written
func (c *Producer) ProduceBytes(msg []byte) error {
	return c.writer.WriteMessages(context.Background(), kafka.Message{Value: msg})
}
