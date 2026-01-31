// Package config implements the command-line interface for BackItUp.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDB    MongoDBConfig
	PostgreSQL PostgreSQLConfig
	MySQL      MySQLConfig

	BackupDir    string
	Compression  bool
	SlackWebhook string
}

type MongoDBConfig struct {
	URI string
}

type PostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		MongoDB: MongoDBConfig{
			URI: getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
		},
		PostgreSQL: PostgreSQLConfig{
			Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
			Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
			User:     getEnvOrDefault("POSTGRES_USER", "postgres"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
		},
		MySQL: MySQLConfig{
			Host:     getEnvOrDefault("MYSQL_HOST", "localhost"),
			Port:     getEnvOrDefault("MYSQL_PORT", "3306"),
			User:     getEnvOrDefault("MYSQL_USER", "root"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Database: os.Getenv("MYSQL_DB"),
		},
		BackupDir:    getEnvOrDefault("BACKUP_DIR", "./backups"),
		Compression:  getEnvOrDefault("COMPRESSION", "true") == "true",
		SlackWebhook: os.Getenv("SLACK_WEBHOOK_URL"),
	}

	return cfg, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
