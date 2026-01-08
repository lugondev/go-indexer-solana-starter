# Go Indexer for Solana Starter Program

A high-performance Solana blockchain event indexer built specifically for the [Starter Program](../starter_program/README.md). This indexer monitors and stores all events emitted by the Starter Program smart contracts into MongoDB or PostgreSQL for easy querying and analysis.

## 🚀 Features

- **Event-Driven Architecture**: Indexes all 20+ event types from Starter Program
- **Real-time Processing**: Polls Solana RPC for new transactions and processes events immediately
- **Multiple Database Support**: MongoDB and PostgreSQL (with migrations)
- **Anchor Event Decoding**: Automatically decodes Anchor framework events with discriminators
- **Concurrent Processing**: Configurable batch size and concurrency for optimal performance
- **Type-Safe Models**: Strongly-typed event models with Borsh serialization
- **Production Ready**: Graceful shutdown, error handling, and comprehensive logging

## 📋 Indexed Events

The indexer tracks all events from the Starter Program:

### Token Events
- `TokensMintedEvent` - SPL token minting
- `TokensTransferredEvent` - Token transfers
- `TokensBurnedEvent` - Token burning
- `DelegateApprovedEvent` - Delegate approval
- `DelegateRevokedEvent` - Delegate revocation
- `TokenAccountClosedEvent` - Account closure
- `TokenAccountFrozenEvent` - Account freeze
- `TokenAccountThawedEvent` - Account thaw

### User Events
- `UserAccountCreatedEvent` - New user registration
- `UserAccountUpdatedEvent` - User profile updates
- `UserAccountClosedEvent` - Account deletion

### Config Events
- `ConfigUpdatedEvent` - Program configuration changes
- `ProgramPausedEvent` - Program pause/unpause

### NFT Events
- `NftCollectionCreatedEvent` - NFT collection creation
- `NftMintedEvent` - NFT minting
- `NftListedEvent` - NFT marketplace listing
- `NftSoldEvent` - NFT sale
- `NftListingCancelledEvent` - Listing cancellation
- `NftOfferCreatedEvent` - Offer creation
- `NftOfferAcceptedEvent` - Offer acceptance

## 🏗️ Project Structure

```
go_indexer/
├── cmd/
│   └── indexer/              # Main application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── decoder/              # Anchor event decoder
│   ├── indexer/              # Core indexer logic
│   ├── models/               # Event models
│   ├── processor/            # Event processor
│   └── repository/           # Database repositories
│       ├── repository.go     # Repository interface
│       ├── mongo.go          # MongoDB implementation
│       └── postgres.go       # PostgreSQL implementation
├── pkg/
│   └── solana/               # Solana RPC client
├── idl/                      # Anchor IDL files
├── tools/                    # Code generation tools
└── .env.example              # Environment variables template
```

## 🛠️ Prerequisites

- **Go 1.24+** (required for latest dependencies)
- **MongoDB 4.4+** or **PostgreSQL 12+**
- **Solana Devnet/Mainnet access** (RPC endpoint)
- **Starter Program deployed** (see [starter_program](../starter_program/))

## 📦 Installation

### 1. Clone and Setup

```bash
cd go_indexer
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` with your settings:

```env
# Solana Configuration
SOLANA_RPC_URL=https://api.devnet.solana.com
SOLANA_WS_URL=wss://api.devnet.solana.com

# Program IDs (from your deployed programs)
STARTER_PROGRAM_ID=gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC
COUNTER_PROGRAM_ID=CounzVsCGF4VzNkAwePKC9mXr6YWiFYF4kLW6YdV8Cc

# Indexer Settings
START_SLOT=0                  # Set to current slot to index from now
POLL_INTERVAL_MS=5000         # Poll every 5 seconds
BATCH_SIZE=20                 # Process 20 transactions per batch
MAX_CONCURRENCY=5             # 5 concurrent workers

# Database (choose one)
DATABASE_TYPE=mongodb
DATABASE_URL=mongodb://localhost:27017
DATABASE_NAME=solana_indexer

# Or PostgreSQL
# DATABASE_TYPE=postgres
# DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
# DATABASE_NAME=solana_indexer

# Server
SERVER_PORT=8080
LOG_LEVEL=info
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Setup Database

#### MongoDB (Recommended)

```bash
# Install MongoDB
brew install mongodb-community@7.0  # macOS
# or
sudo apt-get install mongodb         # Linux

# Start MongoDB
brew services start mongodb-community@7.0  # macOS
# or
sudo systemctl start mongodb               # Linux

# Verify connection
mongosh --eval "db.version()"
```

#### PostgreSQL (Alternative)

```bash
# Install PostgreSQL
brew install postgresql@15    # macOS
# or
sudo apt-get install postgresql  # Linux

# Start PostgreSQL
brew services start postgresql@15  # macOS
# or
sudo systemctl start postgresql    # Linux

# Create database
createdb solana_indexer

# Run migrations (automatic on first start)
```

## 🚀 Usage

### Running the Indexer

```bash
# Build
go build -o indexer cmd/indexer/main.go

# Run
./indexer

# Or run directly
go run cmd/indexer/main.go
```

### Output Example

```
2026/01/08 15:30:45 starting indexer for program gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC from slot 0
2026/01/08 15:30:50 processing 15 signatures
2026/01/08 15:30:51 processed event TokensMintedEvent at slot 123456
2026/01/08 15:30:51 processed event UserAccountCreatedEvent at slot 123457
2026/01/08 15:30:51 processed event NftMintedEvent at slot 123458
```

## 📊 Querying Events

### MongoDB Queries

```javascript
// Connect to MongoDB
mongosh solana_indexer

// Find all token mint events
db.events.find({ event_type: "TokensMintedEvent" })

// Find events by user
db.events.find({ "user": "USER_PUBKEY_HERE" })

// Find recent events (last 24 hours)
db.events.find({
  block_time: { $gte: new Date(Date.now() - 24*60*60*1000) }
}).sort({ block_time: -1 })

// Aggregate events by type
db.events.aggregate([
  { $group: { _id: "$event_type", count: { $sum: 1 } } }
])

// Find NFT sales above 1 SOL
db.events.find({
  event_type: "NftSoldEvent",
  price: { $gte: 1000000000 }  // 1 SOL in lamports
})
```

### PostgreSQL Queries

```sql
-- Connect to PostgreSQL
psql -d solana_indexer

-- Find all token mint events
SELECT * FROM events 
WHERE event_type = 'TokensMintedEvent';

-- Find events by signature
SELECT * FROM events 
WHERE signature = 'SIGNATURE_HERE';

-- Find recent events (last 24 hours)
SELECT * FROM events 
WHERE block_time > NOW() - INTERVAL '24 hours'
ORDER BY block_time DESC;

-- Count events by type
SELECT event_type, COUNT(*) 
FROM events 
GROUP BY event_type;

-- Find NFT mints with metadata
SELECT 
  signature,
  event_data->>'nft_mint' as nft_mint,
  event_data->>'name' as name,
  event_data->>'uri' as uri,
  block_time
FROM events 
WHERE event_type = 'NftMintedEvent'
ORDER BY block_time DESC;
```

## 🏗️ Architecture

### Event Processing Flow

```
┌─────────────┐
│   Solana    │
│  Blockchain │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────┐
│         RPC Client (pkg/solana)         │
│  - GetSignaturesForAddress()            │
│  - GetTransaction()                     │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│     Decoder (internal/decoder)          │
│  - ParseProgramData() from logs         │
│  - DecodeEvent() with discriminators    │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│   Event Processor (internal/processor)  │
│  - ProcessEvent()                       │
│  - Route to specific handlers           │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│   Repository (internal/repository)      │
│  - SaveEvent()                          │
│  - MongoDB or PostgreSQL                │
└─────────────────────────────────────────┘
```

### Anchor Event Decoding

Events are decoded using Anchor's discriminator system:

1. Extract "Program data:" logs from transaction
2. Base64 decode the data
3. Read first 8 bytes as discriminator (SHA256 hash of "event:EventName")
4. Match discriminator to event type
5. Borsh deserialize remaining bytes into typed struct
6. Store in database with metadata (signature, slot, timestamp)

## 🔧 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/decoder/
go test ./internal/processor/
```

### Code Quality

```bash
# Format code
gofmt -w .

# Run linters
golangci-lint run

# Check for issues
go vet ./...
```

### Adding New Event Types

1. Add event struct to `internal/models/events.go`
2. Add event type constant
3. Add discriminator to `internal/decoder/anchor_decoder.go`
4. Implement decoder function
5. Add handler in `internal/processor/event_processor.go`

## 🐳 Docker Deployment

```bash
# Build image
docker build -t go-indexer .

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f indexer
```

## 📈 Performance Tuning

### Configuration Tips

- **POLL_INTERVAL_MS**: Lower = more real-time, higher = less RPC calls
- **BATCH_SIZE**: Higher = fewer RPC calls but more memory
- **MAX_CONCURRENCY**: Match to your CPU cores (usually 4-8)

### MongoDB Optimization

```javascript
// Create compound indexes for common queries
db.events.createIndex({ event_type: 1, block_time: -1 })
db.events.createIndex({ "user": 1, block_time: -1 })
db.events.createIndex({ program_id: 1, slot: -1 })
```

### PostgreSQL Optimization

```sql
-- Create indexes for common queries
CREATE INDEX idx_events_type_time ON events(event_type, block_time DESC);
CREATE INDEX idx_events_jsonb_user ON events USING GIN ((event_data->'user'));
```

## 🔍 Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### Metrics (TODO)

Future support for Prometheus metrics:
- Events processed per second
- RPC request latency
- Database write latency
- Event types distribution

## 🐛 Troubleshooting

### Common Issues

**"Program data not found in logs"**
- Check if transactions have events
- Verify program ID is correct
- Ensure you're indexing the right transactions

**"Failed to decode event"**
- Check IDL matches deployed program version
- Verify discriminator calculation
- Check Borsh serialization format

**"Too many RPC requests"**
- Increase POLL_INTERVAL_MS
- Reduce BATCH_SIZE
- Use rate-limited RPC endpoint

**"Database connection failed"**
- Verify DATABASE_URL is correct
- Check database is running
- Verify credentials

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## 📚 Resources

- [Starter Program Documentation](../starter_program/README.md)
- [Anchor Framework](https://www.anchor-lang.com/)
- [Solana RPC API](https://solana.com/docs/rpc)
- [go-carbon Framework](https://github.com/lugondev/go-carbon)
