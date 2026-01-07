package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Solana RPC configuration
	SolanaRPCURL string
	SolanaWSURL  string

	// Indexer configuration
	StartSlot      uint64
	PollInterval   time.Duration
	BatchSize      int
	MaxConcurrency int

	// Database configuration
	DatabaseURL string

	// Server configuration
	ServerPort int

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		SolanaRPCURL:   getEnvOrDefault("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com"),
		SolanaWSURL:    getEnvOrDefault("SOLANA_WS_URL", "wss://api.mainnet-beta.solana.com"),
		StartSlot:      uint64(getEnvIntOrDefault("START_SLOT", 0)),
		PollInterval:   time.Duration(getEnvIntOrDefault("POLL_INTERVAL_MS", 1000)) * time.Millisecond,
		BatchSize:      getEnvIntOrDefault("BATCH_SIZE", 10),
		MaxConcurrency: getEnvIntOrDefault("MAX_CONCURRENCY", 5),
		DatabaseURL:    getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/solana_indexer?sslmode=disable"),
		ServerPort:     getEnvIntOrDefault("SERVER_PORT", 8080),
		LogLevel:       getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SolanaRPCURL == "" {
		return fmt.Errorf("SOLANA_RPC_URL is required")
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
