package main

import (
	"github.com/darth-raijin/gafka-binance/internal/config"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/kafka"
	"github.com/darth-raijin/gafka-binance/internal/migrations"
	log "github.com/sirupsen/logrus"
)

func main() {
	appConfig := config.LoadConfig("local")

	db, err := database.NewDatabase(database.PostgresConfig{
		Host:     appConfig.Database.Host,
		User:     appConfig.Database.User,
		Password: appConfig.Database.Password,
		DBName:   appConfig.Database.DBName,
		Port:     appConfig.Database.Port,
		SSLMode:  appConfig.Database.SSLMode,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = migrations.Migrate(db.Connection)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}

	preparedConsumers := initializeConsumers()

}

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
