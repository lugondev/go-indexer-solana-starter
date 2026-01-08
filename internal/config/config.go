package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseType string

const (
	DatabaseTypeMongo    DatabaseType = "mongodb"
	DatabaseTypePostgres DatabaseType = "postgres"
)

type Config struct {
	SolanaRPCURL string
	SolanaWSURL  string

	StarterProgramID string
	CounterProgramID string

	StartSlot      uint64
	PollInterval   time.Duration
	BatchSize      int
	MaxConcurrency int

	DatabaseType DatabaseType
	DatabaseURL  string
	DatabaseName string

	ServerPort int
	LogLevel   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		SolanaRPCURL:     getEnvOrDefault("SOLANA_RPC_URL", "https://api.devnet.solana.com"),
		SolanaWSURL:      getEnvOrDefault("SOLANA_WS_URL", "wss://api.devnet.solana.com"),
		StarterProgramID: getEnvOrDefault("STARTER_PROGRAM_ID", "gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC"),
		CounterProgramID: getEnvOrDefault("COUNTER_PROGRAM_ID", "CounzVsCGF4VzNkAwePKC9mXr6YWiFYF4kLW6YdV8Cc"),
		StartSlot:        uint64(getEnvIntOrDefault("START_SLOT", 0)),
		PollInterval:     time.Duration(getEnvIntOrDefault("POLL_INTERVAL_MS", 1000)) * time.Millisecond,
		BatchSize:        getEnvIntOrDefault("BATCH_SIZE", 10),
		MaxConcurrency:   getEnvIntOrDefault("MAX_CONCURRENCY", 5),
		DatabaseType:     DatabaseType(getEnvOrDefault("DATABASE_TYPE", "mongodb")),
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", "mongodb://localhost:27017"),
		DatabaseName:     getEnvOrDefault("DATABASE_NAME", "solana_indexer"),
		ServerPort:       getEnvIntOrDefault("SERVER_PORT", 8080),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.SolanaRPCURL == "" {
		return fmt.Errorf("SOLANA_RPC_URL is required")
	}
	if c.StarterProgramID == "" {
		return fmt.Errorf("STARTER_PROGRAM_ID is required")
	}
	if c.BatchSize <= 0 {
		return fmt.Errorf("BATCH_SIZE must be positive")
	}
	if c.MaxConcurrency <= 0 {
		return fmt.Errorf("MAX_CONCURRENCY must be positive")
	}
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}
	if c.DatabaseType != DatabaseTypeMongo && c.DatabaseType != DatabaseTypePostgres {
		return fmt.Errorf("DATABASE_TYPE must be 'mongodb' or 'postgres'")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.DatabaseName == "" {
		return fmt.Errorf("DATABASE_NAME is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
			return intVal
		}
	}
	return defaultValue
}
