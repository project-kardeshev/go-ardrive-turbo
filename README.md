# Go Turbo SDK 🚀

A Go SDK for interacting with the Turbo Upload and Payment Services, providing functionality for uploading data to Arweave through the Turbo infrastructure.

## 🛠️ Installation

```bash
go get github.com/project-kardeshev/go-ardrive-turbo
```

## 📋 Available APIs

### Unauthenticated Client

```go
client := turbo.Unauthenticated(nil)
```

#### `GetUploadCosts(ctx, request)`
Estimate the cost in Winston Credits for uploading data of specific sizes.

```go
costs, err := client.GetUploadCosts(ctx, &turbo.UploadCostsRequest{
    Bytes: []int64{1024, 2048, 4096}, // 1KB, 2KB, 4KB
})
```

#### `UploadSignedDataItem(ctx, request)`
Upload a pre-signed data item with progress tracking.

```go
result, err := client.UploadSignedDataItem(ctx, &turbo.SignedDataItemUploadRequest{
    DataItemStreamFactory: func() (io.ReadCloser, error) {
        return os.Open("signed-data-item.bin")
    },
    DataItemSizeFactory: func() int64 {
        return fileSize
    },
    Events: &turbo.UploadEvents{
        OnUploadProgress: func(event turbo.ProgressEvent) {
            fmt.Printf("Upload: %d/%d bytes\n", event.ProcessedBytes, event.TotalBytes)
        },
    },
})
```

### Authenticated Client

```go
client, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
    PrivateKey: jwkMap,           // Arweave JWK
    Token:      turbo.TokenTypeArweave,
})

// Or with Ethereum
client, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
    PrivateKey: "0x1234...",      // Ethereum private key hex
    Token:      turbo.TokenTypeEthereum,
})
```

#### `GetBalance(ctx)`
Get the credit balance of the authenticated wallet in Winston Credits.

```go
balance, err := client.GetBalance(ctx)
fmt.Printf("Balance: %s winc\n", balance.Winc.String())
```

#### `Upload(ctx, request)`
Sign and upload data with automatic signing and progress tracking.

```go
result, err := client.Upload(ctx, &turbo.UploadRequest{
    Data: []byte("Hello, Turbo!"),
    Tags: []turbo.Tag{
        {Name: "Content-Type", Value: "text/plain"},
        {Name: "App-Name", Value: "my-app"},
    },
    Events: &turbo.UploadEvents{
        OnProgress: func(event turbo.ProgressEvent) {
            fmt.Printf("%s: %d/%d bytes\n", event.Step, event.ProcessedBytes, event.TotalBytes)
        },
        OnSigningSuccess: func() {
            fmt.Println("Data signed successfully!")
        },
        OnUploadSuccess: func() {
            fmt.Println("Upload completed!")
        },
    },
})
```

#### `GetSigner()`
Get the signer instance associated with the authenticated client.

```go
signer := client.GetSigner()
address, err := signer.GetNativeAddress()
```

### Signers

#### Arweave Signer

```go
// From JWK
signer, err := turbo.NewArweaveSigner(jwkMap)

// From keyfile
signer, err := turbo.NewArweaveSignerFromKeyfile("path/to/wallet.json")
```

#### Ethereum Signer

```go
signer, err := turbo.NewEthereumSigner("0x1234567890abcdef...")
```

## 🚀 Quick Start Examples

### Get Upload Costs

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    turbo "github.com/project-kardeshev/go-ardrive-turbo"
)

func main() {
    ctx := context.Background()
    client := turbo.Unauthenticated(nil)
    
    costs, err := client.GetUploadCosts(ctx, &turbo.UploadCostsRequest{
        Bytes: []int64{1024, 2048, 4096},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for i, cost := range *costs {
        fmt.Printf("Cost for %d bytes: %s winc\n", 
            []int64{1024, 2048, 4096}[i], cost.Winc.String())
    }
}
```

### Upload Data with Arweave Wallet

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    
    turbo "github.com/project-kardeshev/go-ardrive-turbo"
)

func main() {
    ctx := context.Background()
    
    // Load your Arweave JWK
    walletBytes, err := os.ReadFile("wallet.json")
    if err != nil {
        log.Fatal(err)
    }
    
    var jwk map[string]interface{}
    if err := json.Unmarshal(walletBytes, &jwk); err != nil {
        log.Fatal(err)
    }
    
    // Create authenticated client
    client, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
        PrivateKey: jwk,
        Token:      turbo.TokenTypeArweave,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Check balance
    balance, err := client.GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %s winc\n", balance.Winc.String())
    
    // Upload data
    result, err := client.Upload(ctx, &turbo.UploadRequest{
        Data: []byte("Hello from Go Turbo SDK!"),
        Tags: []turbo.Tag{
            {Name: "Content-Type", Value: "text/plain"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Upload successful! ID: %s\n", result.ID)
}
```

### Upload with Ethereum Wallet

```go
client, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
    PrivateKey: "0x1234567890abcdef...", // Your Ethereum private key
    Token:      turbo.TokenTypeEthereum,
})

result, err := client.Upload(ctx, &turbo.UploadRequest{
    Data: []byte("Hello from Ethereum!"),
    Tags: []turbo.Tag{
        {Name: "Content-Type", Value: "text/plain"},
        {Name: "Wallet-Type", Value: "ethereum"},
    },
})
```

## 📊 Types Reference

### Core Types

```go
type Balance struct {
    Winc Winston `json:"winc"`
}

type UploadResult struct {
    ID                  string   `json:"id"`
    Owner               string   `json:"owner"`
    DataCaches          []string `json:"dataCaches,omitempty"`
    FastFinalityIndexes []string `json:"fastFinalityIndexes,omitempty"`
}

type UploadRequest struct {
    Data    []byte        // Data to upload
    Tags    []Tag         // Arweave tags
    Target  string        // Optional target address
    Anchor  string        // Optional anchor
    Events  *UploadEvents // Progress callbacks
    Context context.Context // Request context
}

type Tag struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}
```

### Event Callbacks

```go
type UploadEvents struct {
    // Overall progress (includes signing + upload)
    OnProgress func(ProgressEvent)
    OnError    func(ErrorEvent)
    OnSuccess  func()
    
    // Signing-specific events
    OnSigningProgress func(ProgressEvent)
    OnSigningError    func(error)
    OnSigningSuccess  func()
    
    // Upload-specific events
    OnUploadProgress func(ProgressEvent)
    OnUploadError    func(error)
    OnUploadSuccess  func()
}

type ProgressEvent struct {
    TotalBytes     int64  `json:"totalBytes"`
    ProcessedBytes int64  `json:"processedBytes"`
    Step           string `json:"step"` // "signing" or "upload"
}
```

### Configuration

```go
// Default configuration (production)
config := turbo.DefaultConfig()
// Uses: https://upload.ardrive.io and https://payment.ardrive.io

// Development configuration
config := turbo.DevConfig()
// Uses: https://upload.ardrive.dev and https://payment.ardrive.dev

// Custom configuration
client := turbo.Unauthenticated(&turbo.UnauthenticatedOptions{
    Config: &turbo.Config{
        UploadURL:  "https://custom-upload.example.com",
        PaymentURL: "https://custom-payment.example.com",
        Token:      turbo.TokenTypeArweave,
    },
})
```

## 🔗 Supported Token Types

- `turbo.TokenTypeArweave` - Arweave (AR)
- `turbo.TokenTypeEthereum` - Ethereum (ETH)
- `turbo.TokenTypeSolana` - Solana (SOL)
- `turbo.TokenTypePolygon` - Polygon (POL/MATIC)
- `turbo.TokenTypeKyve` - KYVE
- `turbo.TokenTypeBaseEth` - Base Ethereum
- `turbo.TokenTypeArio` - AR.IO

## 🧪 Testing

```bash
# Run tests
go test ./tests/

# Run with verbose output
go test -v ./tests/

# Run example
go run examples/basic/main.go
```

The test suite includes a `test_wallet.json` file for comprehensive testing with a real Arweave wallet.

## 📁 Project Structure

```
go-ardrive-turbo/
├── pkg/
│   ├── types/          # Core types and interfaces
│   ├── signers/        # Wallet signers (Arweave, Ethereum)
│   └── turbo/          # Client implementations
├── examples/           # Usage examples
├── tests/             # Unit tests and test wallet
└── turbo.go           # Main package exports
```

## 🔧 Dependencies

- **goar** - Arweave wallet operations and data item signing
- **go-ethereum** - Ethereum wallet operations and signing
- **Standard library** - HTTP client, JSON, crypto operations

## 🤝 Contributing

This SDK is built with Go best practices:

- Clean interfaces and modular design
- Comprehensive error handling
- Context support for cancellation/timeouts
- Strong typing throughout the API
- Extensive test coverage

## 📄 License

[Add your license here]

## 🔗 Links

- [Original TypeScript SDK](https://github.com/ardriveapp/turbo-sdk)
- [Turbo Documentation](https://docs.ardrive.io/)
- [Arweave Documentation](https://docs.arweave.org/)