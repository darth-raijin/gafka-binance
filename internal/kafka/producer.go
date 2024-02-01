package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
	"go.uber.org/zap"
	"log"
)

type KafkaProducer struct {
	producer      *kafka.Producer
	topic         Topic
	binanceClient binance.BinanceInterface
	logger        *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(kafkaCfg KafkaConfig, topic Topic, logger *zap.Logger, binanceClient binance.BinanceInterface) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaCfg.Brokers[0],
		"client.id":         "gafka-binance",
		"acks":              "all",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	return &KafkaProducer{
		producer:      producer,
		topic:         topic,
		binanceClient: binanceClient,
		logger:        logger,
	}
}

// SendMessage sends a message to the KafkaProducer's topic
func (kp *KafkaProducer) SendMessage(message []byte) error {
	deliveryChan := make(chan kafka.Event)

	topic := kp.topic.String()

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		kp.logger.Error("Failed to produce message", zap.Error(err))
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

func (kp *KafkaProducer) Start() {
	// Create a channel for receiving trade data
	tradesChan := make(chan binance.GetAgregateTradeStreamsResponse)
	// Start receiving trade streams in a separate goroutine
	go func() {
		err := kp.binanceClient.GetAggregateTradeStreams("ethusdc", tradesChan)
		if err != nil {
			kp.logger.Error("Failed to start aggregate trade streams", zap.Error(err))
		}
	}()

	for trade := range tradesChan {
		message, err := json.Marshal(trade)
		if err != nil {
			kp.logger.Error("Failed to serialize trade data", zap.Error(err))
			continue
		}

		if err := kp.SendMessage(message); err != nil {
			kp.logger.Error("Failed to send message", zap.Error(err))
		}

		// Optionally, implement a mechanism to handle clean shutdown, e.g., listening for a quit signal
	}
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
