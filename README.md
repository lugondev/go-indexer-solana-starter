# Go Indexer Solana Starter

A high-performance Solana blockchain indexer written in Go.

## Features

- 🚀 High-performance concurrent block processing
- 🔄 Automatic retry and error handling
- 📊 Real-time slot tracking
- 🛡️ Graceful shutdown handling
- 🧪 Comprehensive test coverage
- 🐳 Docker support

## Project Structure

```
.
├── cmd/
│   └── indexer/          # Main application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── indexer/          # Core indexer logic
│   ├── repository/       # Data access layer
│   └── handler/          # HTTP handlers
├── pkg/
│   └── solana/           # Solana client library
├── api/                  # API definitions
├── configs/              # Configuration files
└── docs/                 # Documentation
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, for data persistence)
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/lugondev/go-indexer-solana-starter.git
cd go-indexer-solana-starter
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment variables:
```bash
cp .env.example .env
```

4. Update `.env` with your configuration

## Usage

### Running Locally

```bash
# Build the binary
make build

# Run the indexer
make run

# Or run directly with go
go run cmd/indexer/main.go
```

### Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## Configuration

Configure the indexer using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SOLANA_RPC_URL` | Solana RPC endpoint | `https://api.mainnet-beta.solana.com` |
| `SOLANA_WS_URL` | Solana WebSocket endpoint | `wss://api.mainnet-beta.solana.com` |
| `START_SLOT` | Starting slot number | `0` |
| `POLL_INTERVAL_MS` | Polling interval in milliseconds | `1000` |
| `BATCH_SIZE` | Number of blocks per batch | `10` |
| `MAX_CONCURRENCY` | Maximum concurrent workers | `5` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://localhost:5432/solana_indexer` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level | `info` |

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover
```

### Code Quality

```bash
# Format code
make fmt

# Run linters
make lint
```

### Installing Development Tools

```bash
make install-tools
```

## Architecture

### Indexer Flow

1. **Initialization**: Load configuration and initialize components
2. **Block Processing**: Poll for new blocks at configured intervals
3. **Concurrent Processing**: Process multiple blocks concurrently
4. **Error Handling**: Retry failed operations with exponential backoff
5. **Graceful Shutdown**: Handle shutdown signals and cleanup resources

### Concurrency Model

- Uses goroutines for parallel block processing
- Channel-based communication for coordination
- Context-based cancellation for graceful shutdown
- Mutex protection for shared state

## Testing

The project includes comprehensive tests:

- Unit tests for individual components
- Integration tests for end-to-end flows
- Table-driven tests for better coverage
- Race condition detection with `-race` flag

## Performance Considerations

- Configurable batch size for optimal throughput
- Connection pooling for database operations
- Efficient memory usage with proper resource cleanup
- Monitoring and profiling support with pprof

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linters
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: https://github.com/lugondev/go-indexer-solana-starter/issues
- Documentation: [docs/](./docs/)

## Roadmap

- [ ] WebSocket support for real-time updates
- [ ] Multiple database backends
- [ ] Metrics and monitoring
- [ ] REST API for querying indexed data
- [ ] GraphQL support
- [ ] Program-specific indexing
- [ ] Token account tracking
