# Phase 2 Todo List - Go Turbo SDK Complete Implementation

This phase implements all remaining features from the TypeScript SDK to create a complete Go port.

## Unauthenticated Client APIs
- [ ] **GetSupportedCurrencies()** - List of currencies supported by Turbo Payment Service
- [ ] **GetSupportedCountries()** - List of countries supported by top-up workflow
- [ ] **GetFiatToAR()** - Current fiat to AR conversion rates from pricing oracles
- [ ] **GetFiatRates()** - Current fiat rates for 1 GiB including fees and adjustments
- [ ] **GetWincForFiat()** - Winston Credits amount for provided fiat currency
- [ ] **GetWincForToken()** - Winston Credits amount for provided token amount
- [ ] **GetFiatEstimateForBytes()** - Fiat price estimate for uploading specified bytes
- [ ] **GetTokenPriceForBytes()** - Token price estimate for uploading specified bytes

## Payment & Top-up APIs
- [ ] **CreateCheckoutSession()** - Create Stripe checkout session for fiat top-ups
- [ ] **TopUpWithTokens()** - Top up wallet with crypto tokens (AR, ETH, SOL, etc.)
- [ ] **SubmitFundTransaction()** - Submit existing funding transaction for processing
- [ ] Support for promo codes in authenticated checkout sessions

## Advanced Upload Features
- [ ] **UploadFile()** - Upload files with stream support and progress tracking
- [ ] **UploadFolder()** - Upload folders with manifest generation and concurrent uploads
- [ ] File stream factory pattern for efficient memory usage
- [ ] Manifest generation with custom index/fallback files
- [ ] Concurrent upload limits and failure handling
- [ ] Content-Type detection and custom tags support

## Credit Sharing System
- [ ] **ShareCredits()** - Share credits with other wallet addresses
- [ ] **RevokeCredits()** - Revoke all shared credits for an address
- [ ] **GetCreditShareApprovals()** - List given/received credit approvals
- [ ] Credit approval expiration handling
- [ ] Paid-by functionality for using shared credits

## Multi-Token Support
- [ ] **Arweave (AR)** - Full support via goar integration
- [ ] **Ethereum (ETH)** - Full support via goethers integration  
- [ ] **Solana (SOL)** - Signer integration and token operations
- [ ] **Polygon (POL/MATIC)** - Support for Polygon network
- [ ] **KYVE** - KYVE network token support with mnemonic handling
- [ ] **Base ETH** - Ethereum on Base network support
- [ ] **AR.IO (ARIO)** - AR.IO network token support

## Event System & Progress Tracking
- [ ] **Progress Events** - Overall, signing, and upload progress tracking
- [ ] **Error Events** - Comprehensive error handling and reporting
- [ ] **Success Events** - Success callbacks for different stages
- [ ] **Event Interfaces** - Clean callback interfaces for all event types
- [ ] **Abort Signals** - Context-based cancellation support

## CLI Tool
- [ ] **balance** - Get wallet balance in Turbo Credits
- [ ] **top-up** - Create Stripe checkout for fiat top-ups
- [ ] **crypto-fund** - Fund wallet with crypto tokens
- [ ] **upload-file** - Upload single files with options
- [ ] **upload-folder** - Upload folders with manifest generation
- [ ] **price** - Get price estimates for various inputs
- [ ] **fiat-estimate** - Get fiat estimates for byte counts
- [ ] **token-price** - Get token prices for byte counts
- [ ] **share-credits** - Share credits with other addresses
- [ ] **revoke-credits** - Revoke shared credits
- [ ] **list-shares** - List credit share approvals

## Advanced Features
- [ ] **Context Support** - Proper Go context usage for cancellation/timeouts
- [ ] **Retry Logic** - Intelligent retry mechanisms for network operations
- [ ] **Rate Limiting** - Built-in rate limiting for API calls
- [ ] **Logging** - Structured logging with configurable levels
- [ ] **Configuration** - Environment-based configuration (dev/prod endpoints)
- [ ] **Connection Pooling** - Efficient HTTP client with connection reuse

## Testing & Documentation
- [ ] **Unit Tests** - Comprehensive unit test coverage
- [ ] **Integration Tests** - Tests against live Turbo services
- [ ] **Benchmark Tests** - Performance benchmarks for upload operations
- [ ] **Example Code** - Working examples for all major features
- [ ] **API Documentation** - Complete godoc documentation
- [ ] **README** - Comprehensive usage guide and examples
- [ ] **Migration Guide** - Guide for TypeScript SDK users

## Quality & Performance
- [ ] **Memory Efficiency** - Streaming uploads without loading full files
- [ ] **Concurrent Safety** - Thread-safe operations where needed
- [ ] **Error Handling** - Comprehensive error types and handling
- [ ] **Input Validation** - Proper validation of all inputs
- [ ] **Security** - Secure handling of private keys and sensitive data

## Success Criteria
✅ Feature parity with TypeScript SDK
✅ Idiomatic Go code following Go best practices
✅ Comprehensive test coverage (>90%)
✅ Complete documentation and examples
✅ CLI tool with all commands functional
✅ Performance benchmarks showing efficient operation
