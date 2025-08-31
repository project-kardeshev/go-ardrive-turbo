package turbo

import (
	"context"
	"io"
	"net/http"
	"strings"
)

// MockHTTPClient implements HTTPClient for testing
type MockHTTPClient struct {
	GetFunc        func(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	PostFunc       func(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	PaymentURL     string
	UploadURL      string
	Responses      map[string]*http.Response
	RequestHistory []MockRequest
}

// MockRequest tracks requests made to the mock client
type MockRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		PaymentURL:     "https://mock-payment.test",
		UploadURL:      "https://mock-upload.test",
		Responses:      make(map[string]*http.Response),
		RequestHistory: make([]MockRequest, 0),
	}
}

func (m *MockHTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	m.RequestHistory = append(m.RequestHistory, MockRequest{
		Method:  "GET",
		URL:     url,
		Headers: headers,
	})

	if m.GetFunc != nil {
		return m.GetFunc(ctx, url, headers)
	}

	if resp, exists := m.Responses[url]; exists {
		return resp, nil
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"status":"ok"}`)),
	}, nil
}

func (m *MockHTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	bodyBytes := []byte{}
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	m.RequestHistory = append(m.RequestHistory, MockRequest{
		Method:  "POST",
		URL:     url,
		Headers: headers,
		Body:    string(bodyBytes),
	})

	if m.PostFunc != nil {
		return m.PostFunc(ctx, url, strings.NewReader(string(bodyBytes)), headers)
	}

	if resp, exists := m.Responses[url]; exists {
		return resp, nil
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"id":"test-upload-id","owner":"test-owner"}`)),
	}, nil
}

func (m *MockHTTPClient) GetPaymentURL() string {
	return m.PaymentURL
}

func (m *MockHTTPClient) GetUploadURL() string {
	return m.UploadURL
}

// SetResponse sets a mock response for a specific URL
func (m *MockHTTPClient) SetResponse(url string, response *http.Response) {
	m.Responses[url] = response
}

// GetLastRequest returns the last request made to the mock client
func (m *MockHTTPClient) GetLastRequest() *MockRequest {
	if len(m.RequestHistory) == 0 {
		return nil
	}
	return &m.RequestHistory[len(m.RequestHistory)-1]
}

// GetRequestCount returns the number of requests made
func (m *MockHTTPClient) GetRequestCount() int {
	return len(m.RequestHistory)
}

// ClearHistory clears the request history
func (m *MockHTTPClient) ClearHistory() {
	m.RequestHistory = make([]MockRequest, 0)
}
