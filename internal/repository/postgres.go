package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lugondev/go-indexer-solana-starter/internal/models"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(connString string) (*PostgresRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PostgresRepository{
		pool: pool,
	}, nil
}

func (r *PostgresRepository) SaveEvent(ctx context.Context, event interface{}) error {
	return fmt.Errorf("postgres repository not fully implemented yet")
}

func (r *PostgresRepository) GetEventsByTimeRange(ctx context.Context, from, to time.Time) ([]models.BaseEvent, error) {
	return nil, fmt.Errorf("postgres repository not fully implemented yet")
}

func (r *PostgresRepository) GetEventsByType(ctx context.Context, eventType models.EventType, limit int) ([]interface{}, error) {
	return nil, fmt.Errorf("postgres repository not fully implemented yet")
}

func (r *PostgresRepository) GetEventBySignature(ctx context.Context, signature string) (interface{}, error) {
	return nil, fmt.Errorf("postgres repository not fully implemented yet")
}

func (r *PostgresRepository) Close(ctx context.Context) error {
	r.pool.Close()
	return nil
}

func (r *PostgresRepository) CreateSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		event_type VARCHAR(100) NOT NULL,
		signature VARCHAR(255) UNIQUE NOT NULL,
		slot BIGINT NOT NULL,
		block_time TIMESTAMP NOT NULL,
		program_id VARCHAR(44) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		raw_data JSONB,
		event_data JSONB NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
	CREATE INDEX IF NOT EXISTS idx_events_block_time ON events(block_time DESC);
	CREATE INDEX IF NOT EXISTS idx_events_slot ON events(slot DESC);
	CREATE INDEX IF NOT EXISTS idx_events_program_id ON events(program_id);
	`

	_, err := r.pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}
