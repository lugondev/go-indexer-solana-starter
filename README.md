# Go Indexer for Solana Programs

A high-performance Solana blockchain event indexer built for both [Starter Program](../starter_program/README.md) and [Counter Program](../starter_program/programs/counter_program/). This indexer monitors and stores all events emitted by the programs into MongoDB or PostgreSQL for easy querying and analysis.

## 🚀 Features

- **Multi-Program Support**: Indexes both Starter Program (Anchor events) and Counter Program (log-based events)
- **Event-Driven Architecture**: Indexes 20+ event types from Starter Program + 6 event types from Counter Program
- **Real-time Processing**: Polls Solana RPC for new transactions and processes events immediately
- **Multiple Database Support**: MongoDB and PostgreSQL (with migrations)
- **Dual Decoding Strategy**: Anchor discriminator-based decoding + log parsing
- **Concurrent Processing**: Configurable batch size and concurrency for optimal performance
- **Type-Safe Models**: Strongly-typed event models with proper serialization
- **Production Ready**: Graceful shutdown, error handling, and comprehensive logging

## 📋 Indexed Events

### Starter Program Events

The indexer tracks all Anchor events from the Starter Program:

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

### Counter Program Events

The indexer parses log messages from Counter Program transactions:

- `CounterInitializedEvent` - Counter account creation
- `CounterIncrementedEvent` - Counter incremented by 1
- `CounterDecrementedEvent` - Counter decremented by 1
- `CounterAddedEvent` - Counter incremented by arbitrary value
- `CounterResetEvent` - Counter reset to 0 (authority only)
- `CounterPaymentReceivedEvent` - Counter incremented with SOL payment

## 🏗️ Project Structure

```
go_indexer/
├── cmd/
│   └── indexer/              # Main application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── decoder/              # Event decoders
│   │   ├── anchor_decoder.go # Starter Program Anchor events
│   │   └── counter_parser.go # Counter Program log parser
│   ├── indexer/              # Core indexer logic (multi-program)
│   ├── models/               # Event models (Starter + Counter)
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
2026/01/08 15:30:45 starting indexer for Starter Program gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC from slot 0
2026/01/08 15:30:45 starting indexer for Counter Program CounzVsCGF4VzNkAwePKC9mXr6YWiFYF4kLW6YdV8Cc from slot 0
2026/01/08 15:30:50 processing 15 starter program signatures
2026/01/08 15:30:51 processed starter event TokensMintedEvent at slot 123456
2026/01/08 15:30:51 processed starter event UserAccountCreatedEvent at slot 123457
2026/01/08 15:30:51 processed starter event NftMintedEvent at slot 123458
2026/01/08 15:30:52 processing 3 counter program signatures
2026/01/08 15:30:52 processed counter event CounterIncrementedEvent at slot 123459
2026/01/08 15:30:52 processed counter event CounterPaymentReceivedEvent at slot 123460
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

// Counter Program - Find all counter increments
db.events.find({ event_type: "CounterIncrementedEvent" })

// Counter Program - Find increments that resulted in value > 100
db.events.find({
  event_type: "CounterIncrementedEvent",
  new_value: { $gt: 100 }
})

// Counter Program - Find all payment events with amount >= 0.01 SOL
db.events.find({
  event_type: "CounterPaymentReceivedEvent",
  payment: { $gte: 10000000 }  // 0.01 SOL in lamports
}).sort({ payment: -1 })

// Counter Program - Track counter value changes over time
db.events.find({
  counter: "COUNTER_PUBKEY_HERE",
  event_type: { $in: ["CounterIncrementedEvent", "CounterDecrementedEvent", "CounterAddedEvent"] }
}).sort({ block_time: 1 })

// Counter Program - Find all reset operations
db.events.find({ event_type: "CounterResetEvent" })
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

-- Counter Program - Find all counter events
SELECT * FROM events
WHERE event_type LIKE 'Counter%'
ORDER BY block_time DESC;

-- Counter Program - Track counter value progression
SELECT 
  signature,
  event_type,
  event_data->>'old_value' as old_value,
  event_data->>'new_value' as new_value,
  block_time
FROM events
WHERE event_data->>'counter' = 'COUNTER_PUBKEY_HERE'
  AND event_type IN ('CounterIncrementedEvent', 'CounterDecrementedEvent', 'CounterAddedEvent')
ORDER BY block_time ASC;

-- Counter Program - Find high-value payments
SELECT 
  signature,
  event_data->>'payer' as payer,
  event_data->>'payment' as payment,
  event_data->>'new_count' as new_count,
  block_time
FROM events
WHERE event_type = 'CounterPaymentReceivedEvent'
  AND (event_data->>'payment')::bigint >= 10000000
ORDER BY (event_data->>'payment')::bigint DESC;
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
       ├─────────────────────┬────────────────────────┐
       │                     │                        │
       ▼                     ▼                        ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│ Starter Program  │  │ Counter Program  │  │  Other Programs  │
│   Transactions   │  │   Transactions   │  │   (Future)       │
└────────┬─────────┘  └────────┬─────────┘  └──────────────────┘
         │                     │
         ▼                     ▼
┌──────────────────┐  ┌──────────────────┐
│ Anchor Decoder   │  │  Log Parser      │
│ - Discriminator  │  │ - Regex Extract  │
│ - Borsh Decode   │  │ - msg!() parsing │
└────────┬─────────┘  └────────┬─────────┘
         │                     │
         └─────────┬───────────┘
                   ▼
         ┌─────────────────────┐
         │  Event Processor    │
         │  - ProcessEvent()   │
         │  - Route handlers   │
         └──────────┬──────────┘
                    ▼
         ┌─────────────────────┐
         │    Repository       │
         │  - SaveEvent()      │
         │  - MongoDB/PG       │
         └─────────────────────┘
```

### Decoding Strategies

#### Starter Program: Anchor Event Decoding

Starter Program emits proper Anchor events using discriminators:

1. Extract "Program data:" logs from transaction
2. Base64 decode the data
3. Read first 8 bytes as discriminator (SHA256 hash of "event:EventName")
4. Match discriminator to event type
5. Borsh deserialize remaining bytes into typed struct
6. Store in database with metadata (signature, slot, timestamp)

**Example Anchor Event:**
```rust
#[event]
pub struct TokensMintedEvent {
    pub mint: Pubkey,
    pub recipient: Pubkey,
    pub amount: u64,
}
```

#### Counter Program: Log Message Parsing

Counter Program uses `msg!()` macro instead of Anchor events, requiring log parsing:

1. Extract "Program log:" messages from transaction logs
2. Use regex patterns to match known log formats:
   - `"Counter initialized"` → CounterInitializedEvent
   - `"Counter incremented to: 42"` → CounterIncrementedEvent
   - `"Added 5 to counter. New value: 47"` → CounterAddedEvent
   - `"Payment of 1000000 lamports received. Counter incremented to: 48"` → CounterPaymentReceivedEvent
3. Extract numeric values and account keys from logs
4. Construct event models with parsed data
5. Store in database

**Example Counter Log:**
```rust
msg!("Counter incremented to: {}", counter.count);
// Parsed as: CounterIncrementedEvent { new_value: counter.count, ... }
```

**Why Two Different Strategies?**
- Starter Program: Production-ready with proper Anchor event system
- Counter Program: Simple demonstration program using basic logging
- Indexer supports both patterns for maximum flexibility

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

// Counter-specific indexes
db.events.createIndex({ counter: 1, block_time: 1 })
db.events.createIndex({ event_type: 1, new_value: -1 })
```

### PostgreSQL Optimization

```sql
-- Create indexes for common queries
CREATE INDEX idx_events_type_time ON events(event_type, block_time DESC);
CREATE INDEX idx_events_jsonb_user ON events USING GIN ((event_data->'user'));

-- Counter-specific indexes
CREATE INDEX idx_events_counter ON events USING GIN ((event_data->'counter'));
CREATE INDEX idx_events_payment ON events((event_data->>'payment')::bigint DESC) 
  WHERE event_type = 'CounterPaymentReceivedEvent';
```

## 🧪 Testing with Local Programs

### Test Starter Program Events

```bash
# Terminal 1: Start local validator
solana-test-validator

# Terminal 2: Deploy and test Starter Program
cd starter_program
anchor build
anchor deploy
anchor test --skip-local-validator

# Terminal 3: Run indexer
cd ../go_indexer
./indexer

# Terminal 4: Query events
mongosh solana_indexer
db.events.find().sort({block_time: -1}).limit(10).pretty()
```

### Test Counter Program Events

```bash
# Terminal 1: Local validator (already running)

# Terminal 2: Run Counter Program tests
cd starter_program
anchor test tests/cross_program.ts --skip-local-validator

# Expected Counter events in indexer logs:
# - CounterInitializedEvent
# - CounterIncrementedEvent
# - CounterAddedEvent
# - CounterPaymentReceivedEvent

# Terminal 3: Verify Counter events in MongoDB
mongosh solana_indexer
db.events.find({ event_type: /Counter/ }).pretty()

# Check counter value progression
db.events.find({ 
  event_type: "CounterIncrementedEvent" 
}).sort({ block_time: 1 }).forEach(e => {
  print(`Slot ${e.slot}: ${e.old_value} -> ${e.new_value}`)
})
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

**"Program data not found in logs"** (Starter Program)
- Check if transactions have Anchor events
- Verify program ID is correct
- Ensure you're indexing the right transactions
- Look for "Program data:" in transaction logs

**"Failed to decode event"** (Starter Program)
- Check IDL matches deployed program version
- Verify discriminator calculation
- Check Borsh serialization format
- Ensure event struct matches on-chain data

**"No Counter events found"** (Counter Program)
- Verify Counter Program transactions are being indexed
- Check logs contain "Program log:" messages
- Ensure regex patterns match actual log format
- Counter Program uses `msg!()` not Anchor events

**"Too many RPC requests"**
- Increase POLL_INTERVAL_MS
- Reduce BATCH_SIZE
- Use rate-limited RPC endpoint
- Consider running separate indexer instances per program

**"Database connection failed"**
- Verify DATABASE_URL is correct
- Check database is running
- Verify credentials
- Test connection: `mongosh $DATABASE_URL` or `psql $DATABASE_URL`

**"Counter values not updating"**
- Check transaction logs contain numeric values
- Verify regex patterns in counter_parser.go
- Enable debug logging to see parsed values
- Ensure accounts array is properly extracted from transaction

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## 📚 Resources

- [Root Project README](../README.md) - Full monorepo overview
- [Starter Program Documentation](../starter_program/README.md) - Program instructions and API
- [Frontend Documentation](../frontend/README.md) - UI components and hooks
- [Cross-Program Invocation Guide](../starter_program/CROSS_PROGRAM.md) - CPI patterns
- [Anchor Framework](https://www.anchor-lang.com/) - Solana program framework
- [Solana RPC API](https://solana.com/docs/rpc) - RPC endpoint documentation
- [MongoDB Documentation](https://www.mongodb.com/docs/) - MongoDB query language
- [PostgreSQL Documentation](https://www.postgresql.org/docs/) - PostgreSQL SQL reference
- [go-carbon Framework](https://github.com/lugondev/go-carbon) - Go application framework

## 📋 Additional Documentation

- **[QUICKSTART.md](./QUICKSTART.md)** - Quick start guide
- **[IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)** - Implementation details
- **[PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)** - Project overview
- **[CHANGELOG.md](./CHANGELOG.md)** - Version history
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guidelines
- **[docs/architecture.md](./docs/architecture.md)** - Architecture deep dive
- **[docs/api.md](./docs/api.md)** - API documentation
- **[docs/deployment.md](./docs/deployment.md)** - Deployment guide
