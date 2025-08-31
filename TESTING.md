# Testing Documentation

## Overview

This project implements comprehensive testing following Go best practices. We maintain high test coverage with a focus on reliability, performance, and maintainability.

## Test Structure

### Unit Tests
Located in each package directory (`*_test.go` files):

- **`pkg/types/types_test.go`** - Tests for type definitions, data structures, and serialization
- **`pkg/signers/types_test.go`** - Tests for signer interfaces, data item creation, and mock objects
- **`pkg/turbo/client_test.go`** - Tests for HTTP client, unauthenticated operations, and JSON parsing
- **`pkg/turbo/authenticated_test.go`** - Tests for authenticated operations, upload workflows, and event handling
- **`pkg/turbo/factory_test.go`** - Tests for client factory methods and configuration management

### Integration Tests
Located in `test/integration_test.go`:

- **Workflow Testing** - End-to-end workflows for both authenticated and unauthenticated clients
- **Component Integration** - Tests for proper interaction between signers, clients, and HTTP layer
- **Event System** - Comprehensive testing of the event callback system
- **Configuration** - Tests for various configuration scenarios and edge cases

### Benchmark Tests
Located in `pkg/turbo/benchmarks_test.go`:

- **Performance Benchmarks** - Client creation, data item processing, signing operations
- **Memory Profiling** - Memory allocation patterns and efficiency
- **Concurrency Tests** - Thread-safety and parallel operation performance
- **Scalability Tests** - Performance with various data sizes (1KB to 1MB)

### Mock Objects
- **`pkg/signers/mocks.go`** - Mock signer implementations for testing
- **`pkg/turbo/mocks.go`** - Mock HTTP client for isolated testing

## Test Coverage

Current coverage metrics:

```
pkg/signers      20.6% coverage
pkg/turbo        77.0% coverage  
pkg/types        100% coverage (no testable statements)
```

### Coverage Goals
- **Critical paths**: 100% coverage for core upload and signing workflows
- **Error handling**: 100% coverage for all error conditions
- **Public APIs**: 100% coverage for all exported functions and methods
- **Edge cases**: Comprehensive coverage for boundary conditions and error scenarios

## Running Tests

### Quick Commands
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

### Using Makefile
```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-integration

# Tests with coverage
make test-coverage

# Benchmarks
make bench

# Full CI pipeline
make ci
```

### Advanced Testing
```bash
# Run specific test
go test -run TestAuthenticatedClientUpload ./pkg/turbo/

# Run specific benchmark
go test -bench BenchmarkSignDataItem ./pkg/turbo/

# Memory profiling
go test -bench . -memprofile mem.prof ./pkg/turbo/

# CPU profiling
go test -bench . -cpuprofile cpu.prof ./pkg/turbo/

# Generate coverage with function details
go test -coverprofile=coverage.out ./pkg/...
go tool cover -func=coverage.out
```

## Test Categories

### 1. Functional Tests
- **Input/Output validation** - Verify correct data processing
- **State management** - Test object state changes and persistence
- **API compliance** - Ensure interface implementations meet contracts

### 2. Error Handling Tests
- **Network errors** - HTTP timeouts, connection failures, invalid responses
- **Validation errors** - Invalid input data, missing required fields
- **Security errors** - Authentication failures, signing errors

### 3. Performance Tests
- **Throughput** - Operations per second under load
- **Latency** - Response times for various operations  
- **Memory efficiency** - Memory allocation patterns and garbage collection
- **Concurrency** - Thread-safety and parallel operation performance

### 4. Integration Tests
- **Component interaction** - Proper integration between layers
- **Configuration** - Various configuration scenarios
- **Event flow** - Event system functionality and callback execution
- **Workflow completeness** - End-to-end operation verification

## Best Practices

### Test Naming
- Use descriptive test names: `TestAuthenticatedClientUploadWithValidData`
- Use table-driven tests for multiple scenarios
- Group related tests with subtests: `t.Run("success case", func(t *testing.T) {...})`

### Test Organization
- One test file per source file
- Group tests by functionality
- Use setup/teardown functions for common initialization
- Keep tests independent and isolated

### Mock Usage
- Use mocks to isolate units under test
- Verify mock interactions when relevant
- Prefer dependency injection for testability
- Reset mocks between tests

### Assertions
- Use clear error messages: `t.Errorf("Expected %v, got %v", expected, actual)`
- Test both success and failure cases
- Verify all important side effects
- Use table-driven tests for multiple inputs

## Continuous Integration

The project uses GitHub Actions for automated testing:

- **Multi-platform testing** - Linux, macOS, Windows
- **Multi-version Go support** - Go 1.21, 1.22, 1.23
- **Automated coverage reporting** - Codecov integration
- **Linting and formatting** - golangci-lint integration
- **Security scanning** - gosec for vulnerability detection

### Local CI Simulation
```bash
# Run full CI pipeline locally
make ci

# Individual CI steps
make fmt-check
make mod-verify
make test-coverage
```

## Test Data Management

### Mock Data
- Use realistic but safe test data
- Avoid hardcoded secrets or real credentials
- Generate test data programmatically when possible
- Keep test data minimal but representative

### Test Fixtures
- Store complex test data in separate files when needed
- Use consistent naming for test fixtures
- Document test data requirements and constraints

## Debugging Tests

### Common Debugging Techniques
```bash
# Run single test with verbose output
go test -v -run TestSpecificTest ./pkg/turbo/

# Add debugging output
t.Logf("Debug info: %+v", debugData)

# Use race detector
go test -race ./...

# Generate test binary for debugging
go test -c ./pkg/turbo/
./turbo.test -test.run TestSpecificTest -test.v
```

### Performance Debugging
```bash
# CPU profiling
go test -bench . -cpuprofile cpu.prof
go tool pprof cpu.prof

# Memory profiling  
go test -bench . -memprofile mem.prof
go tool pprof mem.prof

# Trace execution
go test -trace trace.out
go tool trace trace.out
```

## Contributing to Tests

When adding new functionality:

1. **Write tests first** (TDD approach recommended)
2. **Ensure high coverage** for new code paths
3. **Add integration tests** for new workflows
4. **Update benchmarks** if performance is affected
5. **Document test scenarios** in code comments

### Test Review Checklist
- [ ] Tests cover all public APIs
- [ ] Error conditions are tested
- [ ] Edge cases are handled
- [ ] Tests are deterministic and repeatable
- [ ] Performance implications are considered
- [ ] Integration points are verified
- [ ] Documentation is updated
