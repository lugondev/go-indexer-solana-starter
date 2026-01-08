#!/bin/bash
set -e

echo "=== Go Indexer Configuration Test ==="
echo ""

# Check Go version
echo "✓ Go version:"
go version

# Check dependencies
echo ""
echo "✓ Checking dependencies..."
go mod verify

# Build
echo ""
echo "✓ Building indexer..."
go build -o indexer cmd/indexer/main.go

# Check binary
echo ""
echo "✓ Binary size:"
ls -lh indexer | awk '{print $5}'

# Verify environment variables
echo ""
echo "✓ Required environment variables:"
grep -v "^#" .env.example | grep "=" | cut -d= -f1 | while read var; do
  echo "  - $var"
done

echo ""
echo "=== Configuration Test Passed! ==="
echo ""
echo "Next steps:"
echo "1. Copy .env.example to .env and configure"
echo "2. Start MongoDB: docker-compose up -d mongodb"
echo "3. Run indexer: ./indexer"
