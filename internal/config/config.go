// Package config implements the command-line interface for BackItUp.
package config

import (
	"log"

	"github.com/spf13/viper"
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
	URI  string
	path string
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

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
		} else {
			log.Fatal("Error reading config:", err)
		}
	}
}

func Load() (*Config, error) {
	cfg := &Config{
		MongoDB: MongoDBConfig{
			URI:  getEnvOrDefault("mongodb.uri", "mongodb://localhost:27017"),
			path: getEnvOrDefault("mongodb.path", ""),
		},
		PostgreSQL: PostgreSQLConfig{
			Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
			Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
			User:     getEnvOrDefault("POSTGRES_USER", "postgres"),
			Password: getEnvOrDefault("POSTGRES_PASSWORD", ""),
			Database: getEnvOrDefault("POSTGRES_DB", ""),
		},
		MySQL: MySQLConfig{
			Host:     getEnvOrDefault("MYSQL_HOST", "localhost"),
			Port:     getEnvOrDefault("MYSQL_PORT", "3306"),
			User:     getEnvOrDefault("MYSQL_USER", "root"),
			Password: getEnvOrDefault("MYSQL_PASSWORD", ""),
			Database: getEnvOrDefault("MYSQL_DB", ""),
		},
		BackupDir:    getEnvOrDefault("BACKUP_DIR", "./backups"),
		Compression:  getEnvOrDefault("COMPRESSION", "true") == "true",
		SlackWebhook: getEnvOrDefault("SLACK_WEBHOOK_URL", ""),
	}

	return cfg, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func SetMongodbURI(uri string) {
	viper.Set("mongodb.uri", uri)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMongodbPath(path string) {
	viper.Set("mongodb.path", path)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMySQLHost(host string) {
	viper.Set("MYSQL_HOST", host)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMySQLPort(port string) {
	viper.Set("MYSQL_PORT", port)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMySQLUser(user string) {
	viper.Set("MYSQL_USER", user)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMySQLPassword(password string) {
	viper.Set("MYSQL_PASSWORD", password)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetMySQLDatabase(database string) {
	viper.Set("MYSQL_DB", database)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetPostgreSQLHost(host string) {
	viper.Set("POSTGRES_HOST", host)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetPostgreSQLPort(port string) {
	viper.Set("POSTGRES_PORT", port)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetPostgreSQLUser(user string) {
	viper.Set("POSTGRES_USER", user)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetPostgreSQLPassword(password string) {
	viper.Set("POSTGRES_PASSWORD", password)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}

func SetPostgreSQLDatabase(database string) {
	viper.Set("POSTGRES_DB", database)

	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatal("Failed to save config:", err)
		}
	}
}
