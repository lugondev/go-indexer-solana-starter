package indexer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lugondev/go-indexer-solana-starter/internal/config"
)

// Indexer represents the main indexer service
type Indexer struct {
	cfg          *config.Config
	currentSlot  uint64
	mu           sync.RWMutex
	isRunning    bool
	shutdownOnce sync.Once
}

// New creates a new Indexer instance
func New(cfg *config.Config) (*Indexer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Indexer{
		cfg:         cfg,
		currentSlot: cfg.StartSlot,
		isRunning:   false,
	}, nil
}

// Start begins the indexing process
func (i *Indexer) Start(ctx context.Context) error {
	i.mu.Lock()
	if i.isRunning {
		i.mu.Unlock()
		return fmt.Errorf("indexer is already running")
	}
	i.isRunning = true
	i.mu.Unlock()

	log.Printf("starting indexer from slot %d", i.currentSlot)

	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("indexer context cancelled")
			return ctx.Err()
		case <-ticker.C:
			if err := i.processBlocks(ctx); err != nil {
				log.Printf("error processing blocks: %v", err)
				// Continue processing despite errors
			}
		}
	}
}

// Shutdown gracefully stops the indexer
func (i *Indexer) Shutdown(ctx context.Context) error {
	var shutdownErr error
	i.shutdownOnce.Do(func() {
		i.mu.Lock()
		defer i.mu.Unlock()

		if !i.isRunning {
			return
		}

		log.Println("shutting down indexer...")
		i.isRunning = false

		// Add cleanup logic here
		// For example: close database connections, flush buffers, etc.
	})
	return shutdownErr
}

// processBlocks processes a batch of blocks
func (i *Indexer) processBlocks(ctx context.Context) error {
	i.mu.RLock()
	currentSlot := i.currentSlot
	batchSize := i.cfg.BatchSize
	i.mu.RUnlock()

	// TODO: Implement actual block processing logic
	log.Printf("processing blocks from slot %d (batch size: %d)", currentSlot, batchSize)

	// Simulate processing
	time.Sleep(100 * time.Millisecond)

	// Update current slot
	i.mu.Lock()
	i.currentSlot += uint64(batchSize)
	i.mu.Unlock()

	return nil
}

// GetCurrentSlot returns the current slot being processed
func (i *Indexer) GetCurrentSlot() uint64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.currentSlot
}

// IsRunning returns whether the indexer is currently running
func (i *Indexer) IsRunning() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.isRunning
}
