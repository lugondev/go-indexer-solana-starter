package repository

import (
	"context"
	"time"

	"github.com/lugondev/go-indexer-solana-starter/internal/models"
)

type Repository interface {
	SaveEvent(ctx context.Context, event interface{}) error
	GetEventsByTimeRange(ctx context.Context, from, to time.Time) ([]models.BaseEvent, error)
	GetEventsByType(ctx context.Context, eventType models.EventType, limit int) ([]interface{}, error)
	GetEventBySignature(ctx context.Context, signature string) (interface{}, error)
	Close(ctx context.Context) error
}
