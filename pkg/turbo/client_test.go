package turbo

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func TestNewDefaultHTTPClient(t *testing.T) {
	paymentURL := "https://payment.test"
	uploadURL := "https://upload.test"

	client := NewDefaultHTTPClient(paymentURL, uploadURL)

	if client.GetPaymentURL() != paymentURL {
		t.Errorf("Expected payment URL '%s', got '%s'", paymentURL, client.GetPaymentURL())
	}

	if client.GetUploadURL() != uploadURL {
		t.Errorf("Expected upload URL '%s', got '%s'", uploadURL, client.GetUploadURL())
	}
}

func TestUnauthenticatedClientGetBalance(t *testing.T) {
	mockClient := NewMockHTTPClient()
	client := NewUnauthenticatedClientForTesting(mockClient)

	// Mock successful response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"winc":"1000000000","credits":"1.0","currency":"USD"}`)),
	}
	mockClient.SetResponse("https://mock-payment.test/v1/account/balance/arweave?address=test-address", mockResponse)

	ctx := context.Background()
	balance, err := client.GetBalance(ctx, "test-address")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if balance == nil {
		t.Error("Expected non-nil balance")
	}

	if balance.WinC != "1000000000" {
		t.Errorf("Expected WinC '1000000000', got '%s'", balance.WinC)
	}

	if balance.Credits != "1.0" {
		t.Errorf("Expected Credits '1.0', got '%s'", balance.Credits)
	}

	if balance.Currency != "USD" {
		t.Errorf("Expected Currency 'USD', got '%s'", balance.Currency)
	}

	// Verify request was made correctly
	if mockClient.GetRequestCount() != 1 {
		t.Errorf("Expected 1 request, got %d", mockClient.GetRequestCount())
	}

	lastRequest := mockClient.GetLastRequest()
	if lastRequest.Method != "GET" {
		t.Errorf("Expected GET request, got %s", lastRequest.Method)
	}

	expectedURL := "https://mock-payment.test/v1/account/balance/arweave?address=test-address"
	if lastRequest.URL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, lastRequest.URL)
	}
}

func TestUnauthenticatedClientGetUploadCosts(t *testing.T) {
	mockClient := NewMockHTTPClient()
	client := NewUnauthenticatedClientForTesting(mockClient)

	// Mock successful responses for individual byte counts
	mockResponse1 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"winc":"1000","bytes":1024}`)),
	}
	mockResponse2 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"winc":"1000000","bytes":1048576}`)),
	}
	mockClient.SetResponse("https://mock-payment.test/v1/price/bytes/1024", mockResponse1)
	mockClient.SetResponse("https://mock-payment.test/v1/price/bytes/1048576", mockResponse2)

	ctx := context.Background()
	bytes := []int64{1024, 1048576}
	costs, err := client.GetUploadCosts(ctx, bytes)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(costs) != 2 {
		t.Errorf("Expected 2 costs, got %d", len(costs))
	}

	if costs[0].Winc != "1000" {
		t.Errorf("Expected first cost Winc '1000', got '%s'", costs[0].Winc)
	}

	if costs[0].Bytes != 1024 {
		t.Errorf("Expected first cost Bytes 1024, got %d", costs[0].Bytes)
	}

	if costs[1].Winc != "1000000" {
		t.Errorf("Expected second cost Winc '1000000', got '%s'", costs[1].Winc)
	}

	if costs[1].Bytes != 1048576 {
		t.Errorf("Expected second cost Bytes 1048576, got %d", costs[1].Bytes)
	}

	// Verify two individual requests were made
	if mockClient.GetRequestCount() != 2 {
		t.Errorf("Expected 2 requests, got %d", mockClient.GetRequestCount())
	}
}

func TestUnauthenticatedClientUploadSignedDataItem(t *testing.T) {
	mockClient := NewMockHTTPClient()
	client := NewUnauthenticatedClientForTesting(mockClient)

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
	mockClient.SetResponse("https://mock-upload.test/v1/tx", mockResponse)

	// Track events
	var progressEvents []types.ProgressEvent
	var uploadSuccessCalled bool
	var uploadStartCalled bool

	req := &types.SignedDataItemUploadRequest{
		DataItemStreamFactory: func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("test-data-item")), nil
		},
		DataItemSizeFactory: func() int64 {
			return 14 // len("test-data-item")
		},
		Events: &types.UploadEvents{
			OnProgress: func(event types.ProgressEvent) {
				progressEvents = append(progressEvents, event)
			},
			OnUploadStart: func() {
				uploadStartCalled = true
			},
			OnUploadSuccess: func(result *types.UploadResult) {
				uploadSuccessCalled = true
			},
		},
	}

	ctx := context.Background()
	result, err := client.UploadSignedDataItem(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if result.ID != "test-upload-id" {
		t.Errorf("Expected ID 'test-upload-id', got '%s'", result.ID)
	}

	if result.Owner != "test-owner" {
		t.Errorf("Expected Owner 'test-owner', got '%s'", result.Owner)
	}

	// Verify events were called
	if !uploadStartCalled {
		t.Error("Expected OnUploadStart to be called")
	}

	if !uploadSuccessCalled {
		t.Error("Expected OnUploadSuccess to be called")
	}

	if len(progressEvents) < 2 {
		t.Errorf("Expected at least 2 progress events, got %d", len(progressEvents))
	}

	// Verify request was made correctly
	lastRequest := mockClient.GetLastRequest()
	if lastRequest.Method != "POST" {
		t.Errorf("Expected POST request, got %s", lastRequest.Method)
	}

	if lastRequest.Headers["Content-Type"] != "application/octet-stream" {
		t.Errorf("Expected Content-Type 'application/octet-stream', got '%s'",
			lastRequest.Headers["Content-Type"])
	}

	if lastRequest.Body != "test-data-item" {
		t.Errorf("Expected body 'test-data-item', got '%s'", lastRequest.Body)
	}
}

func TestUnauthenticatedClientUploadSignedDataItemNilRequest(t *testing.T) {
	mockClient := NewMockHTTPClient()
	client := NewUnauthenticatedClientForTesting(mockClient)

	ctx := context.Background()
	_, err := client.UploadSignedDataItem(ctx, nil)

	if err == nil {
		t.Error("Expected error for nil request")
	}

	if !strings.Contains(err.Error(), "upload request is required") {
		t.Errorf("Expected 'upload request is required' error, got '%v'", err)
	}
}

func TestParseJSONError(t *testing.T) {
	// Test HTTP error response
	errorResponse := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`{"error":"Bad Request"}`)),
	}

	var result map[string]interface{}
	err := ParseJSON(errorResponse, &result)

	if err == nil {
		t.Error("Expected error for HTTP 400 response")
	}

	if !strings.Contains(err.Error(), "HTTP 400") {
		t.Errorf("Expected 'HTTP 400' in error, got '%v'", err)
	}
}

func TestParseJSONInvalidJSON(t *testing.T) {
	// Test invalid JSON response
	invalidJSONResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{invalid json`)),
	}

	var result map[string]interface{}
	err := ParseJSON(invalidJSONResponse, &result)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "failed to decode JSON response") {
		t.Errorf("Expected JSON decode error, got '%v'", err)
	}
}

func TestMockHTTPClient(t *testing.T) {
	mock := NewMockHTTPClient()

	// Test default URLs
	if mock.GetPaymentURL() != "https://mock-payment.test" {
		t.Errorf("Expected default payment URL, got '%s'", mock.GetPaymentURL())
	}

	if mock.GetUploadURL() != "https://mock-upload.test" {
		t.Errorf("Expected default upload URL, got '%s'", mock.GetUploadURL())
	}

	// Test request tracking
	ctx := context.Background()
	headers := map[string]string{"Authorization": "Bearer token"}

	mock.Get(ctx, "https://test.com", headers)

	if mock.GetRequestCount() != 1 {
		t.Errorf("Expected 1 request, got %d", mock.GetRequestCount())
	}

	lastRequest := mock.GetLastRequest()
	if lastRequest.Method != "GET" {
		t.Errorf("Expected GET method, got %s", lastRequest.Method)
	}

	if lastRequest.URL != "https://test.com" {
		t.Errorf("Expected URL 'https://test.com', got '%s'", lastRequest.URL)
	}

	if lastRequest.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header, got '%s'", lastRequest.Headers["Authorization"])
	}

	// Test clear history
	mock.ClearHistory()
	if mock.GetRequestCount() != 0 {
		t.Errorf("Expected 0 requests after clear, got %d", mock.GetRequestCount())
	}
}
