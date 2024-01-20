package database

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// Database struct holds the GORM Connection instance.
type Database struct {
	Connection *gorm.DB
}

// Struct for PostgreSQL connection variables.
type PostgresConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Port     string `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
}

// Connection string builder for PostgreSQL.
func (p PostgresConfig) DSN() string {
	return "host=" + p.Host + " user=" + p.User + " password=" + p.Password + " dbname=" + p.DBName + " port=" + p.Port + " sslmode=" + p.SSLMode
}

// NewDatabase establishes a connection with a PostgreSQL database using GORM.
func NewDatabase(cfg PostgresConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build the connection string.
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ping the database to check if the connection is actually open.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{Connection: db}, nil
}
