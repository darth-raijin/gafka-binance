package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    Topic
}

// NewConsumer creates a new Kafka consumer for a given topic
func NewConsumer(kafkaCfg KafkaConfig, topic Topic) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaCfg.Brokers[0],
		"group.id":          fmt.Sprintf("%s-%s", kafkaCfg.ConsumerSettings.GroupID, topic),
		"auto.offset.reset": kafkaCfg.ConsumerSettings.AutoOffsetReset,
	})

	if err != nil {
		log.Fatalf("Failed to create consumer for topic %s: %s\n", topic, err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
	}
}

// Start begins consuming messages on the KafkaConsumer's topic
func (kc *KafkaConsumer) Start(processFunc func([]byte)) {
	err := kc.consumer.SubscribeTopics([]string{string(kc.topic)}, nil)
	if err != nil {
		return
	}
	log.Printf("Consumer started for topic: %s\n", kc.topic)

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err == nil {
			processFunc(msg.Value)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

// Close closes the Kafka consumer
func (kc *KafkaConsumer) Close() {
	err := kc.consumer.Close()
	if err != nil {
		return
	}
}
