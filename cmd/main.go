package main

import (
	"github.com/darth-raijin/gafka-binance/internal/config"
	"github.com/darth-raijin/gafka-binance/internal/database"
	"github.com/darth-raijin/gafka-binance/internal/migrations"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
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

	zap, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		return
	}

	err = migrations.Migrate(db.Connection, zap)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}

	wireModules(db, zap)
}
