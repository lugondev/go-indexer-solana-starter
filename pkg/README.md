# Public Packages

Reusable packages that can be imported by external projects.

## Structure

### solana/
Solana blockchain client library.

**Features:**
- RPC client for interacting with Solana nodes
- Type definitions for Solana data structures
- WebSocket support (planned)

**Usage:**
```go
import "github.com/lugondev/go-indexer-solana-starter/pkg/solana"

client, err := solana.NewClient(
    "https://api.mainnet-beta.solana.com",
    "wss://api.mainnet-beta.solana.com",
)
if err != nil {
    log.Fatal(err)
}

slot, err := client.GetSlot(context.Background())
block, err := client.GetBlock(context.Background(), slot)
```

## Design Principles

Packages in `pkg/` should:
- Have no dependencies on `internal/` packages
- Be well-documented with godoc comments
- Have comprehensive tests
- Follow semantic versioning
- Be stable and backward-compatible

## Adding New Packages

When adding a new package to `pkg/`:
1. Ensure it's truly reusable
2. Write comprehensive documentation
3. Add examples in godoc
4. Achieve >80% test coverage
5. Consider API stability
