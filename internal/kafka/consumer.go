package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
	"github.com/darth-raijin/gafka-binance/internal/models"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    Topic
	logger   *zap.Logger
	db       *database.Database
}

// NewConsumer creates a new Kafka consumer for a given topic
func NewConsumer(kafkaCfg KafkaConfig, topic Topic, logger *zap.Logger, db *database.Database) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaCfg.Brokers[0],
		"group.id":          fmt.Sprintf("%s-%s", kafkaCfg.ConsumerSettings.GroupID, topic),
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		logger.Error("Failed to create consumer", zap.Error(err))
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
		logger:   logger,
		db:       db,
	}
}

// Start begins consuming messages on the KafkaConsumer's topic
func (kc *KafkaConsumer) Start() {
	err := kc.consumer.SubscribeTopics([]string{string(kc.topic)}, nil)
	if err != nil {
		return
	}
	log.Printf("Consumer started for topic: %s\n", kc.topic)

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err == nil {
			kc.logger.Info("Message received", zap.String("topic", string(kc.topic)), zap.String("message", string(msg.Value)))
			kc.processMessage(msg.Value)
		} else {
			kc.logger.Error("Consumer error", zap.Error(err))
		}
	}
}

func (kc *KafkaConsumer) processMessage(msg []byte) {

	// Parse the message into a Trade struct
	// Save the trade to the database
	event := &binance.GetAgregateTradeStreamsResponse{}
	err := json.Unmarshal(msg, event)
	if err != nil {
		kc.logger.Error("Failed to unmarshal message", zap.Error(err))
		return
	}

	parsedPrice, err := strconv.ParseFloat(event.Price, 64)
	if err != nil {
		parsedPrice = 0
	}

	parsedQuantity, err := strconv.ParseFloat(event.Quantity, 64)
	if err != nil {
		parsedQuantity = 0
	}

	parsedEvent := models.Trade{
		Ticker:           event.Symbol,
		AggrerateTradeID: event.TradeID,
		Price:            parsedPrice,
		Quantity:         parsedQuantity,
		BuyerOrderID:     event.BuyerOrderID,
		SellerOrderID:    event.SellerOrderID,
		MarketBuyer:      event.MarketBuyer,
	}

	err = kc.db.Connection.Create(&parsedEvent).Error
	if err != nil {
		kc.logger.Error("Failed to save trade to database", zap.Error(err))
		return
	}
	kc.logger.Info("Trade saved to database", zap.String("ticker", parsedEvent.Ticker), zap.Int("trade_id", parsedEvent.AggrerateTradeID))
}

// Close closes the Kafka consumer
func (kc *KafkaConsumer) Close() {
	err := kc.consumer.Close()
	if err != nil {
		kc.logger.Error("Failed to close consumer", zap.Error(err))
		return
	}
}
