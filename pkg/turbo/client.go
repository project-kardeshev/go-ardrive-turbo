package turbo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
