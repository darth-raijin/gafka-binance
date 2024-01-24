package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
)

type KafkaProducer struct {
	producer      *kafka.Producer
	topic         Topic
	binanceClient *binance.BinanceClient
}

// NewProducer creates a new Kafka producer
func NewProducer(kafkaCfg KafkaConfig, topic Topic) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaCfg.Brokers[0],
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}
}

// SendMessage sends a message to the KafkaProducer's topic
func (kp *KafkaProducer) SendMessage(message []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: (*string)(&kp.topic), Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
