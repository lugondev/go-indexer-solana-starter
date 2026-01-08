# Quick Start Guide

This guide will help you get the Solana Programs indexer running in 5 minutes.

## Prerequisites

- Go 1.24+
- MongoDB or Docker

## Option 1: Local Setup (Recommended for Development)

### 1. Install MongoDB

**macOS:**
```bash
brew tap mongodb/brew
brew install mongodb-community@7.0
brew services start mongodb-community@7.0
```

**Ubuntu/Debian:**
```bash
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
```

### 2. Configure Environment

```bash
cd go_indexer
cp .env.example .env
```

Edit `.env` if needed (defaults are fine for devnet):
- `STARTER_PROGRAM_ID` - Your deployed Starter Program ID
- `COUNTER_PROGRAM_ID` - Your deployed Counter Program ID
- `DATABASE_URL` - MongoDB connection string

### 3. Run the Indexer

```bash
# Install dependencies
go mod download

# Run
go run cmd/indexer/main.go
```

You should see:
```
2026/01/08 15:30:45 starting indexer for Starter Program gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC from slot 0
2026/01/08 15:30:45 starting indexer for Counter Program CounzVsCGF4VzNkAwePKC9mXr6YWiFYF4kLW6YdV8Cc from slot 0
2026/01/08 15:30:50 processing 15 starter program signatures
2026/01/08 15:30:51 processed starter event TokensMintedEvent at slot 123456
2026/01/08 15:30:52 processing 3 counter program signatures
2026/01/08 15:30:52 processed counter event CounterIncrementedEvent at slot 123457
```

## Option 2: Docker (Recommended for Production)

### 1. Configure Environment

```bash
cd go_indexer
cp .env.example .env
# Edit .env with your program ID
```

### 2. Start with Docker Compose

```bash
# Start MongoDB + Indexer
docker-compose up -d

# View logs
docker-compose logs -f indexer

# Stop
docker-compose down
```

## Verify It's Working

### Check MongoDB

```bash
# Connect to MongoDB
mongosh solana_indexer

# Count indexed events
db.events.countDocuments()

# View recent events
db.events.find().sort({block_time: -1}).limit(5).pretty()
```

## Query Examples

### Starter Program Events

```javascript
// Find token mint events
db.events.find({ event_type: "TokensMintedEvent" }).pretty()

// Find events by user
db.events.find({ "user": "YOUR_PUBKEY_HERE" })
```

### Counter Program Events

```javascript
// Find all counter increments
db.events.find({ event_type: "CounterIncrementedEvent" }).pretty()

// Find payment events
db.events.find({ event_type: "CounterPaymentReceivedEvent" }).sort({ payment: -1 })

// Track counter value changes
db.events.find({ 
  counter: "COUNTER_PUBKEY_HERE",
  event_type: /Counter(Incremented|Decremented|Added)/
}).sort({ block_time: 1 })
```

### Statistics

```javascript
// Get event statistics
db.events.aggregate([
  { $group: { _id: "$event_type", count: { $sum: 1 } } },
  { $sort: { count: -1 } }
])
```

## Troubleshooting

### "Failed to connect to database"
- Make sure MongoDB is running: `brew services list` (macOS) or `systemctl status mongod` (Linux)
- Check DATABASE_URL in .env

### "No signatures found"
- Program might not have any transactions yet
- Verify STARTER_PROGRAM_ID and COUNTER_PROGRAM_ID are correct
- Check RPC endpoint is accessible
- Try generating transactions with `anchor test`

### "Failed to decode event"
- Make sure you're using the correct program IDs
- For Starter Program: Check IDL matches deployed version
- For Counter Program: Check log message formats match parser regex

### "Counter events not appearing"
- Counter Program uses log parsing, not Anchor events
- Verify "Program log:" messages are in transaction logs
- Check regex patterns in `internal/decoder/counter_parser.go`
- Enable debug logging to see parsed log messages

## Testing with Local Programs

### Generate Test Transactions

```bash
# Terminal 1: Start local validator
solana-test-validator

# Terminal 2: Deploy and test programs
cd ../starter_program
anchor build
anchor deploy
anchor test --skip-local-validator

# Terminal 3: Run indexer (in another terminal)
cd ../go_indexer
./indexer

# Terminal 4: Watch events appear
mongosh solana_indexer
db.events.find().sort({block_time: -1}).limit(10).pretty()
```

## Next Steps

1. **Generate historical data**: Set `START_SLOT` to an earlier slot in `.env`
2. **Add alerting**: Modify `internal/processor/event_processor.go` to add custom logic
3. **Build API**: Add REST endpoints to query events
4. **Add metrics**: Integrate Prometheus for monitoring

## Support

For issues:
- Check [README.md](./README.md) for full documentation
- Review logs: `docker-compose logs indexer` or check console output
- Verify MongoDB: `mongosh solana_indexer` and run queries
