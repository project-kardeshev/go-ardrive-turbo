package turbo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// unauthenticatedClient implements TurboUnauthenticatedClient as a standalone client
type unauthenticatedClient struct {
	client     *http.Client
	paymentURL string
	uploadURL  string
	token      string
}

// NewUnauthenticatedClient creates a new unauthenticated Turbo client
func NewUnauthenticatedClient(paymentURL, uploadURL string) TurboUnauthenticatedClient {
	return NewUnauthenticatedClientWithToken(paymentURL, uploadURL, "arweave")
}

// NewUnauthenticatedClientWithToken creates a new unauthenticated Turbo client with token type
func NewUnauthenticatedClientWithToken(paymentURL, uploadURL, token string) TurboUnauthenticatedClient {
	return &unauthenticatedClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		paymentURL: paymentURL,
		uploadURL:  uploadURL,
		token:      token,
	}
}

// NewUnauthenticatedClientForTesting creates a new unauthenticated Turbo client with HTTPClient injection for testing
func NewUnauthenticatedClientForTesting(httpClient HTTPClient) TurboUnauthenticatedClient {
	return &testableUnauthenticatedClient{
		httpClient: httpClient,
	}
}

// testableUnauthenticatedClient is a test-friendly implementation that wraps HTTPClient
type testableUnauthenticatedClient struct {
	httpClient HTTPClient
}

// GetBalance implementation for testable client
func (c *testableUnauthenticatedClient) GetBalance(ctx context.Context, address string) (*types.Balance, error) {
	url := fmt.Sprintf("%s/v1/account/balance/arweave?address=%s", c.httpClient.GetPaymentURL(), address)
	resp, err := c.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	var balance types.Balance
	if err := ParseJSON(resp, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// GetUploadCosts implementation for testable client
func (c *testableUnauthenticatedClient) GetUploadCosts(ctx context.Context, bytes []int64) ([]types.UploadCost, error) {
	// Make individual requests for each byte count (matching TypeScript implementation)
	costs := make([]types.UploadCost, len(bytes))
	
	for i, byteCount := range bytes {
		url := fmt.Sprintf("%s/v1/price/bytes/%d", c.httpClient.GetPaymentURL(), byteCount)
		resp, err := c.httpClient.Get(ctx, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get upload cost for byte count %d: %w", byteCount, err)
		}

		var cost types.UploadCost
		if err := ParseJSON(resp, &cost); err != nil {
			return nil, fmt.Errorf("failed to parse response for byte count %d: %w", byteCount, err)
		}
		
		costs[i] = cost
	}

	return costs, nil
}

// UploadSignedDataItem implementation for testable client
func (c *testableUnauthenticatedClient) UploadSignedDataItem(ctx context.Context, req *types.SignedDataItemUploadRequest) (*types.UploadResult, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request is required")
	}

	// Get data stream
	dataStream, err := req.DataItemStreamFactory()
	if err != nil {
		return nil, fmt.Errorf("failed to create data stream: %w", err)
	}
	defer dataStream.Close()

	// Notify upload start
	if req.Events != nil && req.Events.OnUploadStart != nil {
		req.Events.OnUploadStart()
	}
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     req.DataItemSizeFactory(),
			ProcessedBytes: 0,
			Step:           "uploading",
		})
	}

	// Upload the data item
	url := fmt.Sprintf("%s/v1/tx", c.httpClient.GetUploadURL())
	resp, err := c.httpClient.Post(ctx, url, dataStream, map[string]string{
		"Content-Type": "application/octet-stream",
	})
	if err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "uploading"})
		}
		return nil, fmt.Errorf("failed to upload data item: %w", err)
	}

	// Parse the response
	var result types.UploadResult
	if err := ParseJSON(resp, &result); err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "uploading"})
		}
		return nil, err
	}

	// Notify upload success
	if req.Events != nil && req.Events.OnUploadSuccess != nil {
		req.Events.OnUploadSuccess(&result)
	}
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     req.DataItemSizeFactory(),
			ProcessedBytes: req.DataItemSizeFactory(),
			Step:           "uploading",
		})
	}

	return &result, nil
}

// GetBalance returns the credit balance for a given address (unauthenticated version)
func (c *unauthenticatedClient) GetBalance(ctx context.Context, address string) (*types.Balance, error) {
	url := fmt.Sprintf("%s/v1/account/balance/%s?address=%s", c.paymentURL, c.token, address)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Handle 404 responses by returning default balance (matching TypeScript implementation)
	if resp.StatusCode == 404 {
		return &types.Balance{
			WinC:     "0",
			Credits:  "0",
			Currency: "USD",
		}, nil
	}

	var balance types.Balance
	if err := ParseJSON(resp, &balance); err != nil {
		return nil, err
	}

	// If balance is empty, return default balance (matching TypeScript implementation)
	if balance.WinC == "" {
		return &types.Balance{
			WinC:     "0",
			Credits:  "0",
			Currency: "USD",
		}, nil
	}

	return &balance, nil
}

// GetUploadCosts returns the estimated cost in Winston Credits for the provided file sizes
func (c *unauthenticatedClient) GetUploadCosts(ctx context.Context, bytes []int64) ([]types.UploadCost, error) {
	// Make individual requests for each byte count (matching TypeScript implementation)
	costs := make([]types.UploadCost, len(bytes))
	
	for i, byteCount := range bytes {
		url := fmt.Sprintf("%s/v1/price/bytes/%d", c.paymentURL, byteCount)
		
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for byte count %d: %w", byteCount, err)
		}
		
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get upload cost for byte count %d: %w", byteCount, err)
		}

		var cost types.UploadCost
		if err := ParseJSON(resp, &cost); err != nil {
			return nil, fmt.Errorf("failed to parse response for byte count %d: %w", byteCount, err)
		}
		
		costs[i] = cost
	}

	return costs, nil
}

// UploadSignedDataItem uploads a pre-signed data item
func (c *unauthenticatedClient) UploadSignedDataItem(ctx context.Context, req *types.SignedDataItemUploadRequest) (*types.UploadResult, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request is required")
	}

	// Get data stream
	dataStream, err := req.DataItemStreamFactory()
	if err != nil {
		return nil, fmt.Errorf("failed to create data stream: %w", err)
	}
	defer dataStream.Close()

	// Notify upload start
	if req.Events != nil && req.Events.OnUploadStart != nil {
		req.Events.OnUploadStart()
	}
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     req.DataItemSizeFactory(),
			ProcessedBytes: 0,
			Step:           "uploading",
		})
	}

	// Upload the data item
	url := fmt.Sprintf("%s/v1/tx", c.uploadURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, dataStream)
	if err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "uploading"})
		}
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/octet-stream")
	
	resp, err := c.client.Do(httpReq)
	if err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "uploading"})
		}
		return nil, fmt.Errorf("failed to upload data item: %w", err)
	}

	// Parse the response
	var result types.UploadResult
	if err := ParseJSON(resp, &result); err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "uploading"})
		}
		return nil, err
	}

	// Notify upload success
	if req.Events != nil && req.Events.OnUploadSuccess != nil {
		req.Events.OnUploadSuccess(&result)
	}
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     req.DataItemSizeFactory(),
			ProcessedBytes: req.DataItemSizeFactory(),
			Step:           "uploading",
		})
	}

	return &result, nil
}
