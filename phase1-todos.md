# Phase 1 Todo List - Go Turbo SDK Core APIs

This phase focuses on implementing the four most important APIs as requested:
- get balance
- get upload cost  
- upload
- uploadSignedDataItem

## Setup & Infrastructure
- [ ] Set up Go project structure with go.mod, basic directories, and dependencies (goar, goethers)
- [ ] Define core Go types and structs for TurboClient, Balance, UploadResult, etc.
- [ ] Implement TurboFactory with Unauthenticated() and Authenticated() constructors
- [ ] Implement HTTP client with proper error handling and authentication headers
- [ ] Integrate goar and goethers signers for Arweave and Ethereum wallet support

## Core APIs Implementation
- [ ] **GetBalance()** - Get the credit balance of a wallet measured in Winston Credits
- [ ] **GetUploadCosts()** - Get estimated cost in Winston Credits for file sizes
- [ ] **Upload()** - Sign and upload data with progress tracking and events
- [ ] **UploadSignedDataItem()** - Upload pre-signed data items with stream support

## Testing
- [ ] Write unit tests for Phase 1 core functionality
- [ ] Integration tests against Turbo dev environment
- [ ] Example usage documentation

## Dependencies
- `goar` - For Arweave wallet signing and transactions
- `goethers` - For Ethereum wallet signing and transactions
- Standard Go libraries for HTTP, JSON, crypto operations

## Success Criteria
✅ All four core APIs working with proper error handling
✅ Support for both Arweave and Ethereum signers via goar/goethers
✅ Comprehensive test coverage
✅ Clean, idiomatic Go code following Go conventions
