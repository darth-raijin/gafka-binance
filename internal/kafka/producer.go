package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/micvbang/go-helpy/booly"
	"github.com/micvbang/go-helpy/inty"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
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

	fmt.Println("RECEIVED MESSAGE", m.TopicPartition)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func (kp *KafkaProducer) Start() {
	ticker := time.NewTicker(time.Second * 1) // Adjust the interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mockData := binance.GetAgregateTradeStreamsResponse{
				EventType:     "Mock",
				EventTim:      int64(inty.Random() * inty.Random()),
				Symbol:        "BTC",
				TradeID:       inty.Random(),
				Price:         "0.051",
				Quantity:      "1200",
				BuyerOrderID:  inty.Random(),
				SellerOrderID: inty.Random(),
				TradeTime:     int64(inty.Random()),
				MarketBuyer:   booly.Random(),
			}

			// Serialize and send the data as a Kafka message
			message, err := json.Marshal(mockData)
			if err != nil {
				kp.logger.Error("Failed to serialize data", zap.Error(err))
				continue
			}

			if err := kp.SendMessage(message); err != nil {
				kp.logger.Error("Failed to send message", zap.Error(err))
			}
			// You can add more cases here, like a quit signal to stop the producer
		}
	}
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
