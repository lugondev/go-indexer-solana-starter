# Quick Start Guide

This guide will help you get the Solana Starter Program indexer running in 5 minutes.

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
- `STARTER_PROGRAM_ID` - Your deployed program ID
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
2026/01/08 15:30:45 starting indexer for program gARh1g6reuvsAHB7DXqiuYzzyiJeoiJmtmCpV8Y5uWC from slot 0
2026/01/08 15:30:50 processing 15 signatures
2026/01/08 15:30:51 processed event TokensMintedEvent at slot 123456
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

### Find token mint events
```javascript
db.events.find({ event_type: "TokensMintedEvent" }).pretty()
```

### Find events by user
```javascript
db.events.find({ "user": "YOUR_PUBKEY_HERE" })
```

### Get event statistics
```javascript
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
- Verify STARTER_PROGRAM_ID is correct
- Check RPC endpoint is accessible

### "Failed to decode event"
- Make sure you're using the correct program ID
- Check that the IDL matches deployed program version

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
