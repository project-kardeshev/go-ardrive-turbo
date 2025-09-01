package turbo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func TestNewAuthenticatedClient(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)

	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	if client == nil {
		t.Error("Expected non-nil authenticated client")
	}

	if client.GetSigner() != mockSigner {
		t.Error("Expected signer to match")
	}
}

func TestAuthenticatedClientGetBalanceForSigner(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Mock successful balance response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"winc":"2000000000","credits":"2.0","currency":"USD"}`)),
	}
	mockHTTPClient.SetResponse("https://mock-payment.test/v1/account/balance/arweave?address=test-address", mockResponse)

	ctx := context.Background()
	balance, err := client.GetBalanceForSigner(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if balance == nil {
		t.Error("Expected non-nil balance")
	}

	if balance.WinC != "2000000000" {
		t.Errorf("Expected WinC '2000000000', got '%s'", balance.WinC)
	}

	// Verify the correct address was used
	lastRequest := mockHTTPClient.GetLastRequest()
	expectedURL := "https://mock-payment.test/v1/account/balance/arweave?address=test-address"
	if lastRequest.URL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, lastRequest.URL)
	}
}

func TestAuthenticatedClientUpload(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Mock successful upload response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"test-upload-id",
			"owner":"test-owner",
			"dataCaches":["cache1"],
			"fastFinalityIndexes":["index1"],
			"deadlineHeight":1000,
			"block":500,
			"validatorSet":["validator1"],
			"timestamp":1234567890
		}`)),
	}
	mockHTTPClient.SetResponse("https://mock-upload.test/v1/tx", mockResponse)

	// Track events
	var signingSuccessCalled bool
	var uploadSuccessCalled bool
	var progressEvents []types.ProgressEvent

	req := &types.UploadRequest{
		Data: []byte("Hello, Turbo!"),
		Tags: []types.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "App-Name", Value: "go-turbo-test"},
		},
		Target: "target-address",
		Anchor: "anchor-value",
		Events: &types.UploadEvents{
			OnProgress: func(event types.ProgressEvent) {
				progressEvents = append(progressEvents, event)
			},
			OnSigningSuccess: func() {
				signingSuccessCalled = true
			},
			OnUploadSuccess: func(result *types.UploadResult) {
				uploadSuccessCalled = true
			},
		},
	}

	ctx := context.Background()
	result, err := client.Upload(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if result.ID != "test-upload-id" {
		t.Errorf("Expected ID 'test-upload-id', got '%s'", result.ID)
	}

	// Verify events were called
	if !signingSuccessCalled {
		t.Error("Expected OnSigningSuccess to be called")
	}

	if !uploadSuccessCalled {
		t.Error("Expected OnUploadSuccess to be called")
	}

	// Should have progress events for both signing and uploading
	if len(progressEvents) < 2 {
		t.Errorf("Expected at least 2 progress events, got %d", len(progressEvents))
	}

	// Check for signing and uploading steps
	foundSigningStep := false
	foundUploadingStep := false
	for _, event := range progressEvents {
		if event.Step == "signing" {
			foundSigningStep = true
		}
		if event.Step == "uploading" {
			foundUploadingStep = true
		}
	}

	if !foundSigningStep {
		t.Error("Expected progress event with 'signing' step")
	}

	if !foundUploadingStep {
		t.Error("Expected progress event with 'uploading' step")
	}

	// Verify upload request was made
	lastRequest := mockHTTPClient.GetLastRequest()
	if lastRequest.Method != "POST" {
		t.Errorf("Expected POST request, got %s", lastRequest.Method)
	}

	expectedURL := "https://mock-upload.test/v1/tx"
	if lastRequest.URL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, lastRequest.URL)
	}

	if lastRequest.Headers["Content-Type"] != "application/octet-stream" {
		t.Errorf("Expected Content-Type 'application/octet-stream', got '%s'",
			lastRequest.Headers["Content-Type"])
	}
}

func TestAuthenticatedClientUploadNilRequest(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	ctx := context.Background()
	_, err := client.Upload(ctx, nil)

	if err == nil {
		t.Error("Expected error for nil request")
	}

	if !strings.Contains(err.Error(), "upload request is required") {
		t.Errorf("Expected 'upload request is required' error, got '%v'", err)
	}
}

func TestAuthenticatedClientUploadWithDataReader(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Mock successful upload response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(`{
			"id":"test-upload-id",
			"owner":"test-owner"
		}`)),
	}
	mockHTTPClient.SetResponse("https://mock-upload.test/v1/tx", mockResponse)

	// Test with DataReader instead of Data
	dataContent := "Hello from DataReader!"
	req := &types.UploadRequest{
		DataReader: strings.NewReader(dataContent),
		Tags: []types.Tag{
			{Name: "Content-Type", Value: "text/plain"},
		},
	}

	ctx := context.Background()
	result, err := client.Upload(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if result.ID != "test-upload-id" {
		t.Errorf("Expected ID 'test-upload-id', got '%s'", result.ID)
	}
}

func TestAuthenticatedClientUploadNoDataProvided(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Request with neither Data nor DataReader
	req := &types.UploadRequest{
		Tags: []types.Tag{
			{Name: "Content-Type", Value: "text/plain"},
		},
	}

	ctx := context.Background()
	_, err := client.Upload(ctx, req)

	if err == nil {
		t.Error("Expected error when no data is provided")
	}

	if !strings.Contains(err.Error(), "either Data or DataReader must be provided") {
		t.Errorf("Expected data requirement error, got '%v'", err)
	}
}

func TestAuthenticatedClientUploadSigningError(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Set up signer to return an error
	mockSigner.SetSignDataItemError(errors.New("signing failed"))

	// Track error events
	var signingErrorCalled bool
	var errorEventCalled bool

	req := &types.UploadRequest{
		Data: []byte("test data"),
		Events: &types.UploadEvents{
			OnSigningError: func(err error) {
				signingErrorCalled = true
			},
			OnError: func(event types.ErrorEvent) {
				errorEventCalled = true
			},
		},
	}

	ctx := context.Background()
	_, err := client.Upload(ctx, req)

	if err == nil {
		t.Error("Expected error from signing failure")
	}

	if !strings.Contains(err.Error(), "failed to sign data item") {
		t.Errorf("Expected signing error, got '%v'", err)
	}

	// The error events should be called
	if !signingErrorCalled {
		t.Error("Expected OnSigningError to be called")
	}

	if !errorEventCalled {
		t.Error("Expected OnError to be called")
	}
}

func TestAuthenticatedClientInheritsUnauthenticatedMethods(t *testing.T) {
	mockHTTPClient := NewMockHTTPClient()
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	client := NewAuthenticatedClientForTesting(mockHTTPClient, mockSigner)

	// Mock response for GetBalance (unauthenticated method)
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"winc":"1000000000","credits":"1.0","currency":"USD"}`)),
	}
	mockHTTPClient.SetResponse("https://mock-payment.test/v1/account/balance/arweave?address=other-address", mockResponse)

	// Test that authenticated client can use unauthenticated methods
	ctx := context.Background()
	balance, err := client.GetBalance(ctx, "other-address")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if balance == nil {
		t.Error("Expected non-nil balance")
	}

	if balance.WinC != "1000000000" {
		t.Errorf("Expected WinC '1000000000', got '%s'", balance.WinC)
	}
}
