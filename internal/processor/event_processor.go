package processor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/lugondev/go-indexer-solana-starter/internal/models"
	"github.com/lugondev/go-indexer-solana-starter/internal/repository"
)

type EventProcessor struct {
	repo      repository.Repository
	programID solana.PublicKey
}

func NewEventProcessor(repo repository.Repository, programID solana.PublicKey) *EventProcessor {
	return &EventProcessor{
		repo:      repo,
		programID: programID,
	}
}

func (p *EventProcessor) ProcessEvent(ctx context.Context, signature string, slot uint64, blockTime time.Time, eventType models.EventType, eventData interface{}) error {
	baseEvent := models.BaseEvent{
		EventType: eventType,
		Signature: signature,
		Slot:      slot,
		BlockTime: blockTime,
		ProgramID: p.programID,
		CreatedAt: time.Now(),
	}

	switch eventType {
	case models.EventTypeTokensMinted:
		return p.processTokensMinted(ctx, baseEvent, eventData)
	case models.EventTypeTokensTransferred:
		return p.processTokensTransferred(ctx, baseEvent, eventData)
	case models.EventTypeTokensBurned:
		return p.processTokensBurned(ctx, baseEvent, eventData)
	case models.EventTypeUserAccountCreated:
		return p.processUserAccountCreated(ctx, baseEvent, eventData)
	case models.EventTypeUserAccountUpdated:
		return p.processUserAccountUpdated(ctx, baseEvent, eventData)
	case models.EventTypeConfigUpdated:
		return p.processConfigUpdated(ctx, baseEvent, eventData)
	case models.EventTypeNftMinted:
		return p.processNftMinted(ctx, baseEvent, eventData)
	case models.EventTypeCounterInitialized:
		return p.processCounterInitialized(ctx, baseEvent, eventData)
	case models.EventTypeCounterIncremented:
		return p.processCounterIncremented(ctx, baseEvent, eventData)
	case models.EventTypeCounterDecremented:
		return p.processCounterDecremented(ctx, baseEvent, eventData)
	case models.EventTypeCounterAdded:
		return p.processCounterAdded(ctx, baseEvent, eventData)
	case models.EventTypeCounterReset:
		return p.processCounterReset(ctx, baseEvent, eventData)
	case models.EventTypeCounterPaymentReceived:
		return p.processCounterPaymentReceived(ctx, baseEvent, eventData)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func (p *EventProcessor) processTokensMinted(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.TokensMintedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processTokensTransferred(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.TokensTransferredEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processTokensBurned(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.TokensBurnedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processUserAccountCreated(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.UserAccountCreatedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processUserAccountUpdated(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.UserAccountUpdatedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processConfigUpdated(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.ConfigUpdatedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processNftMinted(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.NftMintedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterInitialized(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterInitializedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterIncremented(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterIncrementedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterDecremented(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterDecrementedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterAdded(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterAddedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterReset(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterResetEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) processCounterPaymentReceived(ctx context.Context, base models.BaseEvent, data interface{}) error {
	event := data.(models.CounterPaymentReceivedEvent)
	event.BaseEvent = base
	return p.repo.SaveEvent(ctx, &event)
}

func (p *EventProcessor) GetEventStats(ctx context.Context, from, to time.Time) (map[models.EventType]int64, error) {
	events, err := p.repo.GetEventsByTimeRange(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("get events by time range: %w", err)
	}

	stats := make(map[models.EventType]int64)
	for _, event := range events {
		stats[event.EventType]++
	}

	return stats, nil
}
