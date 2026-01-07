package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("SOLANA_RPC_URL", "https://test.solana.com")
	t.Setenv("START_SLOT", "1000")
	t.Setenv("BATCH_SIZE", "20")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.SolanaRPCURL != "https://test.solana.com" {
		t.Errorf("SolanaRPCURL = %v, want %v", cfg.SolanaRPCURL, "https://test.solana.com")
	}

	if cfg.StartSlot != 1000 {
		t.Errorf("StartSlot = %v, want %v", cfg.StartSlot, 1000)
	}

	if cfg.BatchSize != 20 {
		t.Errorf("BatchSize = %v, want %v", cfg.BatchSize, 20)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
				StartSlot:      0,
				PollInterval:   time.Second,
				BatchSize:      10,
				MaxConcurrency: 5,
				ServerPort:     8080,
			},
			wantErr: false,
		},
		{
			name: "empty RPC URL",
			cfg: &Config{
				SolanaRPCURL:   "",
				BatchSize:      10,
				MaxConcurrency: 5,
				ServerPort:     8080,
			},
			wantErr: true,
		},
		{
			name: "invalid batch size",
			cfg: &Config{
				SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
				BatchSize:      0,
				MaxConcurrency: 5,
				ServerPort:     8080,
			},
			wantErr: true,
		},
		{
			name: "invalid concurrency",
			cfg: &Config{
				SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
				BatchSize:      10,
				MaxConcurrency: -1,
				ServerPort:     8080,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			cfg: &Config{
				SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
				BatchSize:      10,
				MaxConcurrency: 5,
				ServerPort:     70000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
