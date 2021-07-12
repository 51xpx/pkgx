package kafkax

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// ListTopics list topics
func ListTopics(brokers []string) error {
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			return err
		}
		defer conn.Close()

		partitions, err := conn.ReadPartitions()
		if err != nil {
			return err
		}

		m := map[string][]kafka.Partition{}

		for _, p := range partitions {
			m[p.Topic] = append(m[p.Topic], p)
		}
		for k, v := range m {
			fmt.Println(k, v, len(v))
		}
	}

	return nil
}
