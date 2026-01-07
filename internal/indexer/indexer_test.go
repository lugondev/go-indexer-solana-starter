package indexer

import (
	"context"
	"testing"
	"time"

	"github.com/lugondev/go-indexer-solana-starter/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
		{
			name: "valid config",
			cfg: &config.Config{
				SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
				StartSlot:      100,
				PollInterval:   time.Second,
				BatchSize:      10,
				MaxConcurrency: 5,
				ServerPort:     8080,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("New() returned nil indexer")
			}
		})
	}
}

func TestIndexer_GetCurrentSlot(t *testing.T) {
	cfg := &config.Config{
		SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
		StartSlot:      100,
		PollInterval:   time.Second,
		BatchSize:      10,
		MaxConcurrency: 5,
		ServerPort:     8080,
	}

	idx, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create indexer: %v", err)
	}

	if got := idx.GetCurrentSlot(); got != 100 {
		t.Errorf("GetCurrentSlot() = %v, want %v", got, 100)
	}
}

func TestIndexer_StartShutdown(t *testing.T) {
	cfg := &config.Config{
		SolanaRPCURL:   "https://api.mainnet-beta.solana.com",
		StartSlot:      0,
		PollInterval:   50 * time.Millisecond,
		BatchSize:      10,
		MaxConcurrency: 5,
		ServerPort:     8080,
	}

	idx, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create indexer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- idx.Start(ctx)
	}()

	// Wait a bit for indexer to start
	time.Sleep(100 * time.Millisecond)

	if !idx.IsRunning() {
		t.Error("indexer should be running")
	}

	// Cancel context to stop indexer
	cancel()

	// Wait for indexer to stop
	<-errChan

	if err := idx.Shutdown(context.Background()); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}
