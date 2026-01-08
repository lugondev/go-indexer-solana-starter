package indexer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/lugondev/go-indexer-solana-starter/internal/config"
	"github.com/lugondev/go-indexer-solana-starter/internal/decoder"
	"github.com/lugondev/go-indexer-solana-starter/internal/processor"
	"github.com/lugondev/go-indexer-solana-starter/internal/repository"
	solanaClient "github.com/lugondev/go-indexer-solana-starter/pkg/solana"
)

type Indexer struct {
	cfg            *config.Config
	client         *solanaClient.Client
	repo           repository.Repository
	eventProcessor *processor.EventProcessor
	eventDecoder   *decoder.EventDecoder
	programID      solana.PublicKey
	currentSlot    uint64
	lastSignature  *solana.Signature
	mu             sync.RWMutex
	isRunning      bool
	shutdownOnce   sync.Once
}

func New(cfg *config.Config) (*Indexer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	client, err := solanaClient.NewClient(cfg.SolanaRPCURL, cfg.SolanaWSURL)
	if err != nil {
		return nil, fmt.Errorf("create solana client: %w", err)
	}

	programID, err := solana.PublicKeyFromBase58(cfg.StarterProgramID)
	if err != nil {
		return nil, fmt.Errorf("parse program ID: %w", err)
	}

	var repo repository.Repository
	switch cfg.DatabaseType {
	case config.DatabaseTypeMongo:
		repo, err = repository.NewMongoRepository(cfg.DatabaseURL, cfg.DatabaseName)
		if err != nil {
			return nil, fmt.Errorf("create mongo repository: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	eventProcessor := processor.NewEventProcessor(repo, programID)
	eventDecoder := decoder.NewEventDecoder()

	return &Indexer{
		cfg:            cfg,
		client:         client,
		repo:           repo,
		eventProcessor: eventProcessor,
		eventDecoder:   eventDecoder,
		programID:      programID,
		currentSlot:    cfg.StartSlot,
		isRunning:      false,
	}, nil
}

func (i *Indexer) Start(ctx context.Context) error {
	i.mu.Lock()
	if i.isRunning {
		i.mu.Unlock()
		return fmt.Errorf("indexer is already running")
	}
	i.isRunning = true
	i.mu.Unlock()

	log.Printf("starting indexer for program %s from slot %d", i.programID.String(), i.currentSlot)

	if mongoRepo, ok := i.repo.(*repository.MongoRepository); ok {
		if err := mongoRepo.CreateIndexes(ctx); err != nil {
			log.Printf("warning: failed to create indexes: %v", err)
		}
	}

	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("indexer context cancelled")
			return ctx.Err()
		case <-ticker.C:
			if err := i.processSignatures(ctx); err != nil {
				log.Printf("error processing signatures: %v", err)
			}
		}
	}
}

func (i *Indexer) processSignatures(ctx context.Context) error {
	i.mu.RLock()
	programID := i.programID
	lastSig := i.lastSignature
	i.mu.RUnlock()

	sigs, err := i.client.GetSignaturesForAddress(ctx, programID, i.cfg.BatchSize, lastSig, nil)
	if err != nil {
		return fmt.Errorf("get signatures: %w", err)
	}

	if len(sigs) == 0 {
		return nil
	}

	log.Printf("processing %d signatures", len(sigs))

	for _, sig := range sigs {
		if err := i.processTransaction(ctx, sig.Signature); err != nil {
			log.Printf("error processing transaction %s: %v", sig.Signature, err)
			continue
		}
	}

	i.mu.Lock()
	i.lastSignature = &sigs[len(sigs)-1].Signature
	i.mu.Unlock()

	return nil
}

func (i *Indexer) processTransaction(ctx context.Context, signature solana.Signature) error {
	tx, err := i.client.GetTransaction(ctx, signature)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if tx == nil || tx.Meta == nil {
		return nil
	}

	blockTime := time.Unix(int64(tx.BlockTime.Time().Unix()), 0)
	slot := tx.Slot

	logs := tx.Meta.LogMessages
	if len(logs) == 0 {
		return nil
	}

	programDataList := decoder.ParseProgramData(logs)

	for _, data := range programDataList {
		eventType, eventData, err := i.eventDecoder.DecodeEvent(data)
		if err != nil {
			log.Printf("failed to decode event: %v", err)
			continue
		}

		if err := i.eventProcessor.ProcessEvent(ctx, signature.String(), slot, blockTime, eventType, eventData); err != nil {
			log.Printf("failed to process event: %v", err)
			continue
		}

		log.Printf("processed event %s at slot %d", eventType, slot)
	}

	return nil
}

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

		if err := i.repo.Close(ctx); err != nil {
			shutdownErr = fmt.Errorf("close repository: %w", err)
		}
	})
	return shutdownErr
}

func (i *Indexer) GetCurrentSlot() uint64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.currentSlot
}

func (i *Indexer) IsRunning() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.isRunning
}
