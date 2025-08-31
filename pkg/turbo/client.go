package turbo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// HTTPClient represents an HTTP client interface
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	GetPaymentURL() string
	GetUploadURL() string
}

// defaultHTTPClient implements HTTPClient using Go's standard http.Client
type defaultHTTPClient struct {
	client     *http.Client
	paymentURL string
	uploadURL  string
}

// NewDefaultHTTPClient creates a new default HTTP client
func NewDefaultHTTPClient(paymentURL, uploadURL string) HTTPClient {
	return &defaultHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		paymentURL: paymentURL,
		uploadURL:  uploadURL,
	}
}

func (c *defaultHTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.client.Do(req)
}

func (c *defaultHTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.client.Do(req)
}

func (c *defaultHTTPClient) GetPaymentURL() string {
	return c.paymentURL
}

func (c *defaultHTTPClient) GetUploadURL() string {
	return c.uploadURL
}

// unauthenticatedClient implements basic Turbo client functionality
type unauthenticatedClient struct {
	httpClient HTTPClient
}

// NewUnauthenticatedClient creates a new unauthenticated Turbo client
func NewUnauthenticatedClient(httpClient HTTPClient) TurboUnauthenticatedClient {
	return &unauthenticatedClient{
		httpClient: httpClient,
	}
}

// GetBalance returns the credit balance for a given address (unauthenticated version)
func (c *unauthenticatedClient) GetBalance(ctx context.Context, address string) (*types.Balance, error) {
	url := fmt.Sprintf("%s/v1/balance?address=%s", c.httpClient.GetPaymentURL(), address)
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

// GetUploadCosts returns the estimated cost in Winston Credits for the provided file sizes
func (c *unauthenticatedClient) GetUploadCosts(ctx context.Context, bytes []int64) ([]types.UploadCost, error) {
	// Convert bytes slice to comma-separated string
	var bytesStrs []string
	for _, b := range bytes {
		bytesStrs = append(bytesStrs, strconv.FormatInt(b, 10))
	}
	bytesParam := strings.Join(bytesStrs, ",")

	url := fmt.Sprintf("%s/v1/price/bytes/%s", c.httpClient.GetPaymentURL(), bytesParam)
	resp, err := c.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload costs: %w", err)
	}

	var costs []types.UploadCost
	if err := ParseJSON(resp, &costs); err != nil {
		return nil, err
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

// ParseJSON parses JSON response body into the provided interface
func ParseJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}
