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
	"github.com/lugondev/go-indexer-solana-starter/internal/models"
	"github.com/lugondev/go-indexer-solana-starter/internal/processor"
	"github.com/lugondev/go-indexer-solana-starter/internal/repository"
	solanaClient "github.com/lugondev/go-indexer-solana-starter/pkg/solana"
)

type Indexer struct {
	cfg              *config.Config
	client           *solanaClient.Client
	repo             repository.Repository
	starterProcessor *processor.EventProcessor
	counterProcessor *processor.EventProcessor
	eventDecoder     *decoder.EventDecoder
	counterLogParser *decoder.CounterLogParser
	starterProgramID solana.PublicKey
	counterProgramID solana.PublicKey
	currentSlot      uint64
	lastStarterSig   *solana.Signature
	lastCounterSig   *solana.Signature
	mu               sync.RWMutex
	isRunning        bool
	shutdownOnce     sync.Once
}

func New(cfg *config.Config) (*Indexer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	client, err := solanaClient.NewClient(cfg.SolanaRPCURL, cfg.SolanaWSURL)
	if err != nil {
		return nil, fmt.Errorf("create solana client: %w", err)
	}

	starterProgramID, err := solana.PublicKeyFromBase58(cfg.StarterProgramID)
	if err != nil {
		return nil, fmt.Errorf("parse starter program ID: %w", err)
	}

	counterProgramID, err := solana.PublicKeyFromBase58(cfg.CounterProgramID)
	if err != nil {
		return nil, fmt.Errorf("parse counter program ID: %w", err)
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

	starterProcessor := processor.NewEventProcessor(repo, starterProgramID)
	counterProcessor := processor.NewEventProcessor(repo, counterProgramID)
	eventDecoder := decoder.NewEventDecoder()
	counterLogParser := decoder.NewCounterLogParser(counterProgramID)

	return &Indexer{
		cfg:              cfg,
		client:           client,
		repo:             repo,
		starterProcessor: starterProcessor,
		counterProcessor: counterProcessor,
		eventDecoder:     eventDecoder,
		counterLogParser: counterLogParser,
		starterProgramID: starterProgramID,
		counterProgramID: counterProgramID,
		currentSlot:      cfg.StartSlot,
		isRunning:        false,
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

	log.Printf("starting indexer for Starter Program %s from slot %d", i.starterProgramID.String(), i.currentSlot)
	log.Printf("starting indexer for Counter Program %s from slot %d", i.counterProgramID.String(), i.currentSlot)

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
			if err := i.processStarterSignatures(ctx); err != nil {
				log.Printf("error processing starter signatures: %v", err)
			}
			if err := i.processCounterSignatures(ctx); err != nil {
				log.Printf("error processing counter signatures: %v", err)
			}
		}
	}
}

func (i *Indexer) processStarterSignatures(ctx context.Context) error {
	i.mu.RLock()
	programID := i.starterProgramID
	lastSig := i.lastStarterSig
	i.mu.RUnlock()

	sigs, err := i.client.GetSignaturesForAddress(ctx, programID, i.cfg.BatchSize, lastSig, nil)
	if err != nil {
		return fmt.Errorf("get signatures: %w", err)
	}

	if len(sigs) == 0 {
		return nil
	}

	log.Printf("processing %d starter program signatures", len(sigs))

	for _, sig := range sigs {
		if err := i.processStarterTransaction(ctx, sig.Signature); err != nil {
			log.Printf("error processing starter transaction %s: %v", sig.Signature, err)
			continue
		}
	}

	i.mu.Lock()
	i.lastStarterSig = &sigs[len(sigs)-1].Signature
	i.mu.Unlock()

	return nil
}

func (i *Indexer) processCounterSignatures(ctx context.Context) error {
	i.mu.RLock()
	programID := i.counterProgramID
	lastSig := i.lastCounterSig
	i.mu.RUnlock()

	sigs, err := i.client.GetSignaturesForAddress(ctx, programID, i.cfg.BatchSize, lastSig, nil)
	if err != nil {
		return fmt.Errorf("get signatures: %w", err)
	}

	if len(sigs) == 0 {
		return nil
	}

	log.Printf("processing %d counter program signatures", len(sigs))

	for _, sig := range sigs {
		if err := i.processCounterTransaction(ctx, sig.Signature); err != nil {
			log.Printf("error processing counter transaction %s: %v", sig.Signature, err)
			continue
		}
	}

	i.mu.Lock()
	i.lastCounterSig = &sigs[len(sigs)-1].Signature
	i.mu.Unlock()

	return nil
}

func (i *Indexer) processStarterTransaction(ctx context.Context, signature solana.Signature) error {
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

		if err := i.starterProcessor.ProcessEvent(ctx, signature.String(), slot, blockTime, eventType, eventData); err != nil {
			log.Printf("failed to process event: %v", err)
			continue
		}

		log.Printf("processed starter event %s at slot %d", eventType, slot)
	}

	return nil
}

func (i *Indexer) processCounterTransaction(ctx context.Context, signature solana.Signature) error {
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

	var accounts []solana.PublicKey
	if tx.Transaction != nil {
		txObj, err := tx.Transaction.GetTransaction()
		if err == nil {
			accounts = txObj.Message.AccountKeys
		}
	}

	actions, err := i.counterLogParser.ParseLogs(logs, accounts)
	if err != nil {
		return fmt.Errorf("parse counter logs: %w", err)
	}

	for _, action := range actions {
		eventData := i.convertCounterActionToEvent(action)
		if err := i.counterProcessor.ProcessEvent(ctx, signature.String(), slot, blockTime, action.Type, eventData); err != nil {
			log.Printf("failed to process counter event: %v", err)
			continue
		}

		log.Printf("processed counter event %s at slot %d", action.Type, slot)
	}

	return nil
}

func (i *Indexer) convertCounterActionToEvent(action decoder.CounterAction) interface{} {
	switch action.Type {
	case models.EventTypeCounterInitialized:
		authority := solana.PublicKey{}
		if action.Authority != nil {
			authority = *action.Authority
		}
		return models.CounterInitializedEvent{
			Counter:      action.Counter,
			Authority:    authority,
			InitialCount: valueOrDefault(action.NewValue, 0),
		}
	case models.EventTypeCounterIncremented:
		return models.CounterIncrementedEvent{
			Counter:  action.Counter,
			OldValue: valueOrDefault(action.OldValue, 0),
			NewValue: valueOrDefault(action.NewValue, 0),
		}
	case models.EventTypeCounterDecremented:
		return models.CounterDecrementedEvent{
			Counter:  action.Counter,
			OldValue: valueOrDefault(action.OldValue, 0),
			NewValue: valueOrDefault(action.NewValue, 0),
		}
	case models.EventTypeCounterAdded:
		return models.CounterAddedEvent{
			Counter:    action.Counter,
			OldValue:   valueOrDefault(action.OldValue, 0),
			AddedValue: valueOrDefault(action.AddedValue, 0),
			NewValue:   valueOrDefault(action.NewValue, 0),
		}
	case models.EventTypeCounterReset:
		authority := solana.PublicKey{}
		if action.Authority != nil {
			authority = *action.Authority
		}
		return models.CounterResetEvent{
			Counter:   action.Counter,
			Authority: authority,
			OldValue:  valueOrDefault(action.OldValue, 0),
		}
	case models.EventTypeCounterPaymentReceived:
		payer := solana.PublicKey{}
		feeCollector := solana.PublicKey{}
		if action.Payer != nil {
			payer = *action.Payer
		}
		if action.FeeCollector != nil {
			feeCollector = *action.FeeCollector
		}
		return models.CounterPaymentReceivedEvent{
			Counter:      action.Counter,
			Payer:        payer,
			FeeCollector: feeCollector,
			Payment:      valueOrDefault(action.Payment, 0),
			NewCount:     valueOrDefault(action.NewValue, 0),
		}
	default:
		return nil
	}
}

func valueOrDefault(ptr *uint64, defaultValue uint64) uint64 {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
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
