package main

import (
	"github.com/darth-raijin/gafka-binance/internal/config"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/integrations/binance"
	"github.com/darth-raijin/gafka-binance/internal/kafka"
	"go.uber.org/zap"
)

func initializeConsumers() []*kafka.KafkaConsumer {
	consumers := []*kafka.KafkaConsumer{}
	for _, topic := range config.AppConfig.Kafka.Topics {
		switch topic {
		case "Orderbook":
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Orderbook))
		case "Trade":
			consumers = append(consumers, kafka.NewConsumer(config.AppConfig.Kafka, kafka.Trade))
		}
	}
	return consumers
}

func initializeProducers() []*kafka.KafkaProducer {
	producers := []*kafka.KafkaProducer{}
	for _, topic := range config.AppConfig.Kafka.Topics {
		switch topic {
		case "Orderbook":
			producers = append(producers, kafka.NewProducer(config.AppConfig.Kafka, kafka.Orderbook))
		case "Trade":
			producers = append(producers, kafka.NewProducer(config.AppConfig.Kafka, kafka.Trade))
		}
	}
	return producers
}

func wireModules(db *database.Database, logger *zap.Logger) {
	if config.AppConfig.Kafka.Enable {
		preparedConsumers := initializeConsumers()
		preparedProducers := initializeProducers()

		// Start the consumers.
		for _, consumer := range preparedConsumers {
			go consumer.Start(func([]byte) {})
		}

		// Start the producers.
		for _, producer := range preparedProducers {
			go producer.SendMessage([]byte("test"))
		}
	}

	binanceClient := binance.NewBinanceClient(db, config.AppConfig.Binance.Basepaths[0], logger)

	binanceClient.Ping()

}
