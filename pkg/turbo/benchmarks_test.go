package turbo

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// BenchmarkUnauthenticatedClientCreation benchmarks client creation
func BenchmarkUnauthenticatedClientCreation(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := Unauthenticated(config)
		_ = client
	}
}

// BenchmarkAuthenticatedClientCreation benchmarks authenticated client creation
func BenchmarkAuthenticatedClientCreation(b *testing.B) {
	config := DefaultConfig()
	signer := signers.NewMockSigner("test-address", types.TokenTypeArweave)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := Authenticated(config, signer)
		_ = client
	}
}

// BenchmarkDataItemCreation benchmarks data item creation
func BenchmarkDataItemCreation(b *testing.B) {
	data := []byte("benchmark test data")
	tags := []types.Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "benchmark-test"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataItem := signers.CreateDataItem(data, tags, "target", "anchor")
		_ = dataItem
	}
}

// BenchmarkSignDataItem benchmarks data item signing
func BenchmarkSignDataItem(b *testing.B) {
	ctx := context.Background()
	signer := signers.NewMockSigner("test-address", types.TokenTypeArweave)
	data := []byte("benchmark test data for signing")
	tags := []types.Tag{
		{Name: "Content-Type", Value: "text/plain"},
	}
	dataItem := signers.CreateDataItem(data, tags, "", "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bundleItem, err := signer.SignDataItem(ctx, dataItem)
		if err != nil {
			b.Fatalf("Signing failed: %v", err)
		}
		_ = bundleItem
	}
}

// BenchmarkUploadRequestCreation benchmarks upload request creation
func BenchmarkUploadRequestCreation(b *testing.B) {
	data := []byte("benchmark upload data")
	tags := []types.Tag{
		{Name: "Content-Type", Value: "application/octet-stream"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &types.UploadRequest{
			Data: data,
			Tags: tags,
		}
		_ = req
	}
}

// BenchmarkSignedDataItemUploadRequestCreation benchmarks signed upload request creation
func BenchmarkSignedDataItemUploadRequestCreation(b *testing.B) {
	testData := "benchmark signed data item"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &types.SignedDataItemUploadRequest{
			DataItemStreamFactory: func() (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader(testData)), nil
			},
			DataItemSizeFactory: func() int64 {
				return int64(len(testData))
			},
		}
		_ = req
	}
}

// BenchmarkEventHandling benchmarks event system performance
func BenchmarkEventHandling(b *testing.B) {
	events := &types.UploadEvents{
		OnProgress: func(event types.ProgressEvent) {
			// Simulate some work
			_ = event.TotalBytes + event.ProcessedBytes
		},
		OnUploadStart: func() {
			// No-op
		},
		OnUploadSuccess: func(result *types.UploadResult) {
			// Simulate some work
			_ = result.ID
		},
	}

	progressEvent := types.ProgressEvent{
		TotalBytes:     1024,
		ProcessedBytes: 512,
		Step:           "testing",
	}

	uploadResult := &types.UploadResult{
		ID:    "benchmark-id",
		Owner: "benchmark-owner",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if events.OnUploadStart != nil {
			events.OnUploadStart()
		}
		if events.OnProgress != nil {
			events.OnProgress(progressEvent)
		}
		if events.OnUploadSuccess != nil {
			events.OnUploadSuccess(uploadResult)
		}
	}
}

// BenchmarkParseJSON benchmarks JSON parsing performance
func BenchmarkParseJSON(b *testing.B) {
	jsonData := `{"winc":"1000000000","credits":"1.0","currency":"USD"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &mockHTTPResponse{
			statusCode: 200,
			body:       strings.NewReader(jsonData),
		}

		var balance types.Balance
		err := ParseJSON(resp.toHTTPResponse(), &balance)
		if err != nil {
			b.Fatalf("JSON parsing failed: %v", err)
		}
	}
}

// BenchmarkConfigCreation benchmarks configuration creation
func BenchmarkConfigCreation(b *testing.B) {
	b.Run("DefaultConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			config := DefaultConfig()
			_ = config
		}
	})

	b.Run("DevConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			config := DevConfig()
			_ = config
		}
	})

	b.Run("CustomConfig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			config := &TurboConfig{
				PaymentURL: "https://custom-payment.test",
				UploadURL:  "https://custom-upload.test",
			}
			_ = config
		}
	})
}

// BenchmarkLargeDataHandling benchmarks handling of larger data sets
func BenchmarkLargeDataHandling(b *testing.B) {
	sizes := []int{
		1024,    // 1KB
		10240,   // 10KB
		102400,  // 100KB
		1048576, // 1MB
	}

	for _, size := range sizes {
		b.Run(toString(size)+"B", func(b *testing.B) {
			data := make([]byte, size)
			// Fill with some pattern
			for i := range data {
				data[i] = byte(i % 256)
			}

			tags := []types.Tag{
				{Name: "Content-Type", Value: "application/octet-stream"},
				{Name: "Data-Size", Value: toString(size)},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dataItem := signers.CreateDataItem(data, tags, "", "")
				_ = dataItem
			}
		})
	}
}

// Helper functions for benchmarks

type mockHTTPResponse struct {
	statusCode int
	body       io.Reader
}

func (m *mockHTTPResponse) toHTTPResponse() *http.Response {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(m.body),
	}
}

func toString(i int) string {
	switch i {
	case 1024:
		return "1K"
	case 10240:
		return "10K"
	case 102400:
		return "100K"
	case 1048576:
		return "1M"
	default:
		return "unknown"
	}
}

// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	signer := signers.NewMockSigner("test-address", types.TokenTypeArweave)
	data := []byte("concurrent test data")
	tags := []types.Tag{
		{Name: "Content-Type", Value: "text/plain"},
	}
	dataItem := signers.CreateDataItem(data, tags, "", "")

	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		for pb.Next() {
			bundleItem, err := signer.SignDataItem(ctx, dataItem)
			if err != nil {
				b.Errorf("Signing failed: %v", err)
			}
			_ = bundleItem
		}
	})
}
