package turbo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// HTTPClient wraps the standard HTTP client with Turbo-specific functionality
type HTTPClient struct {
	client *http.Client
	config *types.Config
}

// newHTTPClient creates a new HTTP client with appropriate timeouts and configuration
func newHTTPClient(config *types.Config) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

// Request represents an HTTP request to be made
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        interface{}
	BodyReader  io.Reader
	ContentType string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Do executes an HTTP request and returns the response
func (h *HTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	var bodyReader io.Reader

	// Handle different body types
	if req.BodyReader != nil {
		bodyReader = req.BodyReader
	} else if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		if req.ContentType == "" {
			req.ContentType = "application/json"
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	}
	
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	httpResp, err := h.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
	}

	// Check for HTTP errors
	if httpResp.StatusCode >= 400 {
		return response, fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(respBody))
	}

	return response, nil
}

// Get performs a GET request
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return h.Do(ctx, &Request{
		Method:  "GET",
		URL:     url,
		Headers: headers,
	})
}

// Post performs a POST request with JSON body
func (h *HTTPClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return h.Do(ctx, &Request{
		Method:  "POST",
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// PostStream performs a POST request with a streaming body
func (h *HTTPClient) PostStream(ctx context.Context, url string, bodyReader io.Reader, contentType string, headers map[string]string) (*Response, error) {
	return h.Do(ctx, &Request{
		Method:      "POST",
		URL:         url,
		BodyReader:  bodyReader,
		ContentType: contentType,
		Headers:     headers,
	})
}

// GetUploadURL returns the upload service URL
func (h *HTTPClient) GetUploadURL() string {
	return h.config.UploadURL
}

// GetPaymentURL returns the payment service URL
func (h *HTTPClient) GetPaymentURL() string {
	return h.config.PaymentURL
}

// ParseJSON parses JSON response body into the provided interface
func ParseJSON(resp *Response, v interface{}) error {
	if err := json.Unmarshal(resp.Body, v); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return nil
}
