# Internal Packages

Private application code that should not be imported by external projects.

## Structure

### config/
Configuration management for the application. Handles environment variable parsing, validation, and provides type-safe access to configuration values.

**Key Features:**
- Environment-based configuration
- Automatic validation on load
- Sensible defaults
- Type-safe access

### indexer/
Core indexer implementation. Contains the main business logic for processing Solana blocks.

**Key Features:**
- Concurrent block processing
- Graceful shutdown
- Thread-safe state management
- Context-based cancellation

### handler/
HTTP request handlers for the REST API (planned).

**Planned Features:**
- Health check endpoints
- Query endpoints for indexed data
- Metrics endpoints

### repository/
Data access layer for interacting with the database (planned).

**Planned Features:**
- PostgreSQL implementation
- Transaction management
- Query optimization
- Connection pooling

## Import Rules

Code in `internal/` can only be imported by code within this project. This is enforced by the Go compiler.

Example:
```go
// ✅ OK - within same project
import "github.com/lugondev/go-indexer-solana-starter/internal/config"

// ❌ NOT OK - from external project
import "github.com/external/project/internal/config"
```
