package config

import (
	"os"
)

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	URI        string
	Database   string
	Collection string
	Username   string
	Password   string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	config := &DatabaseConfig{
		URI:        getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		Database:   getEnvOrDefault("MONGO_DATABASE", "testdb"),
		Collection: getEnvOrDefault("MONGO_COLLECTION", "users"),
		Username:   os.Getenv("MONGO_USERNAME"),
		Password:   os.Getenv("MONGO_PASSWORD"),
	}

	// If username and password are provided, construct authenticated URI
	if config.Username != "" && config.Password != "" {
		config.URI = "mongodb://" + config.Username + ":" + config.Password + "@localhost:27017"
	}

	return config
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
