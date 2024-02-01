package main

import (
	"github.com/darth-raijin/gafka-binance/internal/config"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
	"github.com/darth-raijin/gafka-binance/internal/kafka"
	"go.uber.org/zap"
)

func initializeConsumers(logger *zap.Logger) []*kafka.KafkaConsumer {
	consumers := []*kafka.KafkaConsumer{}
	for _, topic := range config.AppConfig.Kafka.Topics {
		switch topic {
		case "Orderbook":
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Orderbook, logger))
		case "Trade":
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Trade, logger))
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
		preparedProducers := initializeProducers(logger, binanceClient)

		// Start the producers.
		for _, producer := range preparedProducers {
			go producer.Start()
		}
	}

	binanceClient.Ping()

}
