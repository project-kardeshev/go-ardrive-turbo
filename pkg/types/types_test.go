package types

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestTokenType(t *testing.T) {
	tests := []struct {
		name     string
		token    TokenType
		expected string
	}{
		{
			name:     "Arweave token type",
			token:    TokenTypeArweave,
			expected: "arweave",
		},
		{
			name:     "Ethereum token type",
			token:    TokenTypeEthereum,
			expected: "ethereum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.token) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.token))
			}
		})
	}
}

func TestTag(t *testing.T) {
	tag := Tag{
		Name:  "Content-Type",
		Value: "text/plain",
	}

	if tag.Name != "Content-Type" {
		t.Errorf("Expected tag name 'Content-Type', got '%s'", tag.Name)
	}

	if tag.Value != "text/plain" {
		t.Errorf("Expected tag value 'text/plain', got '%s'", tag.Value)
	}
}

func TestBalance(t *testing.T) {
	balance := Balance{
		WinC:     "1000000000",
		Credits:  "1.0",
		Currency: "USD",
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
}

func TestUploadCost(t *testing.T) {
	cost := UploadCost{
		Winc:  "1000000",
		Bytes: 1024,
		Adjustments: map[string]interface{}{
			"discount": 0.1,
		},
	}

	if cost.Winc != "1000000" {
		t.Errorf("Expected Winc '1000000', got '%s'", cost.Winc)
	}

	if cost.Bytes != 1024 {
		t.Errorf("Expected Bytes 1024, got %d", cost.Bytes)
	}

	if cost.Adjustments["discount"] != 0.1 {
		t.Errorf("Expected discount 0.1, got %v", cost.Adjustments["discount"])
	}
}

func TestProgressEvent(t *testing.T) {
	event := ProgressEvent{
		TotalBytes:     1024,
		ProcessedBytes: 512,
		Step:           "uploading",
	}

	if event.TotalBytes != 1024 {
		t.Errorf("Expected TotalBytes 1024, got %d", event.TotalBytes)
	}

	if event.ProcessedBytes != 512 {
		t.Errorf("Expected ProcessedBytes 512, got %d", event.ProcessedBytes)
	}

	if event.Step != "uploading" {
		t.Errorf("Expected Step 'uploading', got '%s'", event.Step)
	}
}

func TestUploadRequest(t *testing.T) {
	ctx := context.Background()
	reader := strings.NewReader("test data")

	req := UploadRequest{
		Data:       []byte("test data"),
		DataReader: reader,
		Tags: []Tag{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Target:  "target-address",
		Anchor:  "anchor-value",
		Context: ctx,
	}

	if string(req.Data) != "test data" {
		t.Errorf("Expected Data 'test data', got '%s'", string(req.Data))
	}

	if req.DataReader != reader {
		t.Errorf("Expected DataReader to match")
	}

	if len(req.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(req.Tags))
	}

	if req.Tags[0].Name != "Content-Type" {
		t.Errorf("Expected tag name 'Content-Type', got '%s'", req.Tags[0].Name)
	}

	if req.Target != "target-address" {
		t.Errorf("Expected Target 'target-address', got '%s'", req.Target)
	}

	if req.Anchor != "anchor-value" {
		t.Errorf("Expected Anchor 'anchor-value', got '%s'", req.Anchor)
	}

	if req.Context != ctx {
		t.Errorf("Expected Context to match")
	}
}

func TestUploadResult(t *testing.T) {
	result := UploadResult{
		ID:                  "test-id",
		Owner:               "test-owner",
		DataCaches:          []string{"cache1", "cache2"},
		FastFinalityIndexes: []string{"index1", "index2"},
		DeadlineHeight:      1000,
		Block:               500,
		ValidatorSet:        []string{"validator1", "validator2"},
		Timestamp:           1234567890,
	}

	if result.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", result.ID)
	}

	if result.Owner != "test-owner" {
		t.Errorf("Expected Owner 'test-owner', got '%s'", result.Owner)
	}

	if len(result.DataCaches) != 2 {
		t.Errorf("Expected 2 DataCaches, got %d", len(result.DataCaches))
	}

	if len(result.FastFinalityIndexes) != 2 {
		t.Errorf("Expected 2 FastFinalityIndexes, got %d", len(result.FastFinalityIndexes))
	}

	if result.DeadlineHeight != 1000 {
		t.Errorf("Expected DeadlineHeight 1000, got %d", result.DeadlineHeight)
	}

	if result.Block != 500 {
		t.Errorf("Expected Block 500, got %d", result.Block)
	}

	if len(result.ValidatorSet) != 2 {
		t.Errorf("Expected 2 ValidatorSet, got %d", len(result.ValidatorSet))
	}

	if result.Timestamp != 1234567890 {
		t.Errorf("Expected Timestamp 1234567890, got %d", result.Timestamp)
	}
}

func TestSignedDataItemUploadRequest(t *testing.T) {
	testData := "test data item"
	testSize := int64(len(testData))

	req := SignedDataItemUploadRequest{
		DataItemStreamFactory: func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(testData)), nil
		},
		DataItemSizeFactory: func() int64 {
			return testSize
		},
		Context: context.Background(),
	}

	// Test factory functions
	stream, err := req.DataItemStreamFactory()
	if err != nil {
		t.Errorf("Expected no error from DataItemStreamFactory, got %v", err)
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		t.Errorf("Expected no error reading stream, got %v", err)
	}

	if string(data) != testData {
		t.Errorf("Expected data '%s', got '%s'", testData, string(data))
	}

	size := req.DataItemSizeFactory()
	if size != testSize {
		t.Errorf("Expected size %d, got %d", testSize, size)
	}

	if req.Context == nil {
		t.Error("Expected Context to be set")
	}
}
