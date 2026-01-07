package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lugondev/go-indexer-solana-starter/internal/config"
	"github.com/lugondev/go-indexer-solana-starter/internal/indexer"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize indexer
	idx, err := indexer.New(cfg)
	if err != nil {
		log.Fatalf("failed to create indexer: %v", err)
	}

	// Start indexer in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := idx.Start(ctx); err != nil {
			errChan <- fmt.Errorf("indexer error: %w", err)
		}
	}()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		log.Printf("indexer failed: %v", err)
		cancel()
	case sig := <-sigChan:
		log.Printf("received signal %v, shutting down gracefully...", sig)
		cancel()
	}

	// Wait for cleanup
	if err := idx.Shutdown(context.Background()); err != nil {
		log.Printf("error during shutdown: %v", err)
	}

	log.Println("indexer stopped successfully")
}
