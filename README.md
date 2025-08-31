# go-ardrive-turbo

Go SDK for interacting with the Turbo Upload and Payment Services. This is a port of the [@ardrive/turbo-sdk](https://github.com/ardriveapp/turbo-sdk) TypeScript SDK.

## Phase 1 Implementation Status

✅ **Completed Features:**
- `getBalance` - Get credit balance for any address (unauthenticated) or authenticated wallet
- `getUploadCosts` - Get estimated costs for uploading data of various sizes
- `upload` - Sign and upload data to Turbo (authenticated)
- `uploadSignedDataItem` - Upload pre-signed data items (unauthenticated)

## Installation

```bash
go get github.com/project-kardeshev/go-ardrive-turbo
```

## Quick Start

### Unauthenticated Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/project-kardeshev/go-ardrive-turbo/pkg/turbo"
)

func main() {
    ctx := context.Background()
    
    // Create unauthenticated client
    client := turbo.Unauthenticated(turbo.DefaultConfig())
    
    // Get balance for an address
    balance, err := client.GetBalance(ctx, "your-arweave-address")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %s winc\n", balance.WinC)
    
    // Get upload costs
    costs, err := client.GetUploadCosts(ctx, []int64{1024, 1024*1024})
    if err != nil {
        log.Fatal(err)
    }
    for i, cost := range costs {
        fmt.Printf("Cost for upload: %s winc\n", cost.Winc)
    }
}
```

### Authenticated Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
    "github.com/project-kardeshev/go-ardrive-turbo/pkg/turbo"
    "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func main() {
    ctx := context.Background()
    
    // Create signer (example with Arweave)
    signer, err := signers.NewArweaveSignerFromKeyfile("path/to/keyfile.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create authenticated client
    client := turbo.Authenticated(turbo.DefaultConfig(), signer)
    
    // Get balance for authenticated wallet
    balance, err := client.GetBalanceForSigner(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %s winc\n", balance.WinC)
    
    // Upload data
    result, err := client.Upload(ctx, &types.UploadRequest{
        Data: []byte("Hello, Turbo!"),
        Tags: []types.Tag{
            {Name: "Content-Type", Value: "text/plain"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Upload ID: %s\n", result.ID)
}
```

## Architecture

### Core Components

- **`pkg/turbo/`** - Main client implementations and factory
  - `TurboUnauthenticatedClient` - Access to public Turbo services
  - `TurboAuthenticatedClient` - Access to authenticated and public services
  - `TurboFactory` - Factory for creating clients

- **`pkg/signers/`** - Wallet signing implementations
  - `ArweaveSigner` - Arweave wallet support
  - `EthereumSigner` - Ethereum wallet support
  - `Signer` - Common interface for all wallet types

- **`pkg/types/`** - Type definitions and data structures

### Configuration

The SDK supports both production and development environments:

```go
// Production (default)
config := turbo.DefaultConfig()

// Development
config := turbo.DevConfig()

// Custom
config := &turbo.TurboConfig{
    PaymentURL: "https://custom-payment.url",
    UploadURL:  "https://custom-upload.url",
}
```

### Supported Signers

- **Arweave**: JWK-based signing
- **Ethereum**: Private key-based signing

## Examples

See `examples/main.go` for comprehensive usage examples.

## Development

### Build

```bash
go build ./...
# or
make build
```

### Testing

#### Quick Test
```bash
go test ./...
# or
make test
```

#### Comprehensive Testing
```bash
# Run all tests with coverage
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run benchmarks
make bench

# Run with race detection
make test-verbose
```

#### Test Coverage
The project maintains high test coverage:
- **pkg/turbo**: 77.0% coverage
- **pkg/signers**: 20.6% coverage  
- **pkg/types**: 100% coverage (no testable statements)

View detailed coverage:
```bash
make test-coverage
open coverage.html
```

#### Test Structure
- **Unit Tests**: `pkg/*/` - Test individual components in isolation
- **Integration Tests**: `test/` - Test component interactions and workflows  
- **Benchmark Tests**: `pkg/*/*_test.go` - Performance and memory usage tests
- **Mock Objects**: Comprehensive mocks for testing without external dependencies

### Local Development

The project uses a replace directive for local development:

```go
replace github.com/project-kardeshev/go-ardrive-turbo => ./
```

### Development Tools

Use the Makefile for common tasks:

```bash
make help              # Show all available commands
make dev-setup         # Set up development environment
make fmt               # Format code
make lint              # Run linter
make ci                # Run full CI pipeline locally
```

## Roadmap

**Phase 2 (Future):**
- Additional signer types (Solana, etc.)
- More upload options and file handling
- Payment and top-up functionality
- CLI tool
- Comprehensive testing

## Contributing

This project follows Go best practices:
- Code to interfaces
- Prefer type safety over runtime safety
- Prefer composition over inheritance
