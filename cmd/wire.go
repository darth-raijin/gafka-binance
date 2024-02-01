package main

import (
	"github.com/darth-raijin/gafka-binance/internal/config"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
	"github.com/darth-raijin/gafka-binance/internal/kafka"
	"go.uber.org/zap"
)

func initializeConsumers(logger *zap.Logger, db *database.Database) []*kafka.KafkaConsumer {
	consumers := []*kafka.KafkaConsumer{}
	for _, topic := range config.AppConfig.Kafka.Topics {
		switch topic {
		case kafka.Orderbook.String():
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Orderbook, logger, db))
		case kafka.Trade.String():
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Trade, logger, db))
		}
	}
	return consumers
}

func initializeProducers(logger *zap.Logger, binanceClient binance.BinanceInterface) []*kafka.KafkaProducer {
	var producers []*kafka.KafkaProducer
	for _, topic := range config.AppConfig.Kafka.Topics {
		switch topic {
		case "Orderbook":
			producers = append(producers, kafka.NewProducer(config.AppConfig.Kafka, kafka.Orderbook, logger, binanceClient))
		case "Trade":
			producers = append(producers, kafka.NewProducer(config.AppConfig.Kafka, kafka.Trade, logger, binanceClient))
		}
	}
	return producers
}

func wireModules(db *database.Database, logger *zap.Logger) {
	binanceClient := binance.NewBinanceClient(db, config.AppConfig.Binance.Basepaths[0], logger)

	if config.AppConfig.Kafka.Enable {
		// Start the producers.
		for _, producer := range initializeProducers(logger, binanceClient) {
			go producer.Start()
		}

		for _, consumer := range initializeConsumers(logger, db) {
			go consumer.Start()
		}
	}

	binanceClient.Ping()

}
