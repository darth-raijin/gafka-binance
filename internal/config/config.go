package config

import (
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration holds the entire application configuration.
type Configuration struct {
	Database database.PostgresConfig `yaml:"database"`
	Kafka    kafka.KafkaConfig       `yaml:"kafka"`
}

var AppConfig Configuration

func LoadConfig(env string) Configuration {
	// Set the environment default to 'local' if not provided
	environment := env
	if environment == "" {
		environment = "local"
	}

	viper.SetConfigName("config." + environment)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file:", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Error("Error unmarshalling config:", err)
	}

	log.Infof("Initialized application configurations for environment: %s", environment)

	return AppConfig
}
