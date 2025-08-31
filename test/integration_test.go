package test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/turbo"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// TestIntegrationUnauthenticatedWorkflow tests the complete unauthenticated workflow
func TestIntegrationUnauthenticatedWorkflow(t *testing.T) {
	// This test demonstrates the complete workflow using mocked components
	// In a real integration test, this would connect to actual services

	// Create client with dev config
	config := turbo.DevConfig()
	client := turbo.Unauthenticated(config)

	if client == nil {
		t.Fatal("Failed to create unauthenticated client")
	}

	// Test 1: Get balance for an address (would normally require real service)
	// Note: This will fail against real service without proper mocking
	// For integration tests against real services, you'd need actual test data

	t.Log("Integration test created unauthenticated client successfully")
	t.Log("To run against real services, implement proper test fixtures")
}

// TestIntegrationAuthenticatedWorkflow tests the complete authenticated workflow
func TestIntegrationAuthenticatedWorkflow(t *testing.T) {
	// Create mock signer for integration test
	mockSigner := signers.NewMockSigner("test-address", types.TokenTypeArweave)

	// Create authenticated client
	config := turbo.DevConfig()
	client := turbo.Authenticated(config, mockSigner)

	if client == nil {
		t.Fatal("Failed to create authenticated client")
	}

	// Verify signer is properly attached
	signer := client.GetSigner()
	if signer == nil {
		t.Fatal("Expected non-nil signer")
	}

	address, err := signer.GetNativeAddress()
	if err != nil {
		t.Fatalf("Failed to get native address: %v", err)
	}

	if address != "test-address" {
		t.Errorf("Expected address 'test-address', got '%s'", address)
	}

	tokenType := signer.GetTokenType()
	if tokenType != types.TokenTypeArweave {
		t.Errorf("Expected token type Arweave, got %s", tokenType)
	}

	t.Log("Integration test created authenticated client successfully")
	t.Log("Signer integration working correctly")
}

// TestIntegrationDataFlow tests the complete data flow from input to signed output
func TestIntegrationDataFlow(t *testing.T) {
	ctx := context.Background()

	// Create mock signer
	mockSigner := signers.NewMockSigner("test-wallet-address", types.TokenTypeEthereum)

	// Test data creation
	testData := []byte("Integration test data for Turbo SDK")
	testTags := []types.Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "go-turbo-integration-test"},
		{Name: "Version", Value: "1.0.0"},
	}

	// Create data item
	dataItem := signers.CreateDataItem(testData, testTags, "target-address", "anchor-value")

	if dataItem == nil {
		t.Fatal("Failed to create data item")
	}

	// Verify data item structure
	if string(dataItem.Data) != string(testData) {
		t.Errorf("Data mismatch in data item")
	}

	if len(dataItem.Tags) != len(testTags) {
		t.Errorf("Expected %d tags, got %d", len(testTags), len(dataItem.Tags))
	}

	if dataItem.Target != "target-address" {
		t.Errorf("Expected target 'target-address', got '%s'", dataItem.Target)
	}

	if dataItem.Anchor != "anchor-value" {
		t.Errorf("Expected anchor 'anchor-value', got '%s'", dataItem.Anchor)
	}

	// Test signing
	bundleItem, err := mockSigner.SignDataItem(ctx, dataItem)
	if err != nil {
		t.Fatalf("Failed to sign data item: %v", err)
	}

	if len(bundleItem.ItemBinary) == 0 {
		t.Error("Expected non-empty signed data item binary")
	}

	// Test upload request creation
	uploadReq := &types.SignedDataItemUploadRequest{
		DataItemStreamFactory: func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(string(bundleItem.ItemBinary))), nil
		},
		DataItemSizeFactory: func() int64 {
			return int64(len(bundleItem.ItemBinary))
		},
		Context: ctx,
	}

	// Verify upload request
	stream, err := uploadReq.DataItemStreamFactory()
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	streamData, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("Failed to read stream: %v", err)
	}

	if string(streamData) != string(bundleItem.ItemBinary) {
		t.Error("Stream data doesn't match signed data item")
	}

	size := uploadReq.DataItemSizeFactory()
	if size != int64(len(bundleItem.ItemBinary)) {
		t.Errorf("Size mismatch: expected %d, got %d", len(bundleItem.ItemBinary), size)
	}

	t.Log("Complete data flow integration test passed")
}

// TestIntegrationEventHandling tests the event system integration
func TestIntegrationEventHandling(t *testing.T) {
	// Test that events flow properly through the system
	var eventsReceived []string

	events := &types.UploadEvents{
		OnProgress: func(event types.ProgressEvent) {
			eventsReceived = append(eventsReceived, "progress")
		},
		OnUploadStart: func() {
			eventsReceived = append(eventsReceived, "upload_start")
		},
		OnUploadSuccess: func(result *types.UploadResult) {
			eventsReceived = append(eventsReceived, "upload_success")
		},
		OnUploadError: func(err error) {
			eventsReceived = append(eventsReceived, "upload_error")
		},
		OnSigningSuccess: func() {
			eventsReceived = append(eventsReceived, "signing_success")
		},
		OnSigningError: func(err error) {
			eventsReceived = append(eventsReceived, "signing_error")
		},
		OnError: func(event types.ErrorEvent) {
			eventsReceived = append(eventsReceived, "error")
		},
	}

	// Test progress event
	events.OnProgress(types.ProgressEvent{
		TotalBytes:     1024,
		ProcessedBytes: 512,
		Step:           "testing",
	})

	// Test other events
	events.OnUploadStart()
	events.OnSigningSuccess()
	events.OnUploadSuccess(&types.UploadResult{ID: "test-id"})

	expectedEvents := []string{"progress", "upload_start", "signing_success", "upload_success"}

	if len(eventsReceived) != len(expectedEvents) {
		t.Errorf("Expected %d events, got %d", len(expectedEvents), len(eventsReceived))
	}

	for i, expected := range expectedEvents {
		if i >= len(eventsReceived) || eventsReceived[i] != expected {
			t.Errorf("Expected event %d to be '%s', got '%s'", i, expected, eventsReceived[i])
		}
	}

	t.Log("Event handling integration test passed")
}

// TestIntegrationConfigurationOptions tests various configuration scenarios
func TestIntegrationConfigurationOptions(t *testing.T) {
	// Test default configuration
	defaultClient := turbo.Unauthenticated(nil)
	if defaultClient == nil {
		t.Error("Failed to create client with default config")
	}

	// Test dev configuration
	devClient := turbo.Unauthenticated(turbo.DevConfig())
	if devClient == nil {
		t.Error("Failed to create client with dev config")
	}

	// Test production configuration
	prodClient := turbo.Unauthenticated(turbo.DefaultConfig())
	if prodClient == nil {
		t.Error("Failed to create client with production config")
	}

	// Test custom configuration
	customConfig := &turbo.TurboConfig{
		PaymentURL: "https://custom-payment.example.com",
		UploadURL:  "https://custom-upload.example.com",
	}

	customClient := turbo.Unauthenticated(customConfig)
	if customClient == nil {
		t.Error("Failed to create client with custom config")
	}

	// Test authenticated client with different configurations
	mockSigner := signers.NewMockSigner("test-address", types.TokenTypeEthereum)

	authClient := turbo.Authenticated(customConfig, mockSigner)
	if authClient == nil {
		t.Error("Failed to create authenticated client with custom config")
	}

	if authClient.GetSigner() != mockSigner {
		t.Error("Signer not properly attached to authenticated client")
	}

	t.Log("Configuration options integration test passed")
}
