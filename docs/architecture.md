# Architecture

## Overview

The Solana Indexer follows a clean architecture pattern with clear separation of concerns.

## Components

### 1. Main Application (`cmd/indexer`)
- Entry point of the application
- Handles initialization and graceful shutdown
- Manages signal handling (SIGTERM, SIGINT)

### 2. Configuration (`internal/config`)
- Centralized configuration management
- Environment variable parsing
- Configuration validation

### 3. Indexer (`internal/indexer`)
- Core business logic
- Block processing coordination
- State management
- Concurrent processing with goroutines

### 4. Solana Client (`pkg/solana`)
- RPC client for Solana blockchain
- WebSocket support (planned)
- Type definitions for Solana data structures

### 5. Repository (planned)
- Data persistence layer
- Database operations
- Transaction management

### 6. Handler (planned)
- HTTP API endpoints
- Request/response handling
- API documentation

## Data Flow

```
1. Main -> Config Loader
2. Config -> Indexer Initialization
3. Indexer -> Solana Client (fetch blocks)
4. Indexer -> Repository (store data)
5. Handler -> Repository (query data)
```

## Concurrency Model

### Goroutines
- Main indexer loop runs in a separate goroutine
- Block processing can spawn multiple worker goroutines
- Each worker processes blocks independently

### Channels
- Error channel for communicating failures
- Signal channel for graceful shutdown
- Done channel for coordination

### Synchronization
- Mutex for protecting shared state (currentSlot)
- WaitGroup for coordinating worker completion
- Context for cancellation propagation

## Error Handling

1. **Recoverable Errors**: Log and retry
2. **Fatal Errors**: Shutdown gracefully
3. **Context Cancellation**: Clean shutdown

## Future Enhancements

- [ ] WebSocket real-time updates
- [ ] Multiple database backends
- [ ] Metrics and monitoring (Prometheus)
- [ ] Distributed tracing
- [ ] Rate limiting
- [ ] Circuit breaker pattern
