package types

import (
	"context"
	"io"
)

// TokenType represents the type of token/blockchain
type TokenType string

const (
	TokenTypeArweave  TokenType = "arweave"
	TokenTypeEthereum TokenType = "ethereum"
)

// Tag represents a key-value pair for metadata
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Balance represents a wallet's credit balance
type Balance struct {
	WinC     string `json:"winc"`
	Credits  string `json:"credits"`
	Currency string `json:"currency"`
}

// ProgressEvent represents upload progress information
type ProgressEvent struct {
	TotalBytes     int64  `json:"totalBytes"`
	ProcessedBytes int64  `json:"processedBytes"`
	Step           string `json:"step"`
}

// ErrorEvent represents an error that occurred during upload
type ErrorEvent struct {
	Error error  `json:"error"`
	Step  string `json:"step"`
}

// UploadEvents contains callback functions for upload events
type UploadEvents struct {
	OnProgress       func(ProgressEvent)
	OnSigningStart   func()
	OnSigningSuccess func()
	OnSigningError   func(error)
	OnError          func(ErrorEvent)
	OnUploadStart    func()
	OnUploadSuccess  func(*UploadResult)
	OnUploadError    func(error)
}

// UploadRequest represents a request to upload data
type UploadRequest struct {
	Data       []byte          `json:"data,omitempty"`
	DataReader io.Reader       `json:"-"`
	Tags       []Tag           `json:"tags,omitempty"`
	Target     string          `json:"target,omitempty"`
	Anchor     string          `json:"anchor,omitempty"`
	Events     *UploadEvents   `json:"-"`
	Context    context.Context `json:"-"`
}

// UploadResult represents the result of an upload operation
type UploadResult struct {
	ID                  string   `json:"id"`
	Owner               string   `json:"owner"`
	DataCaches          []string `json:"dataCaches"`
	FastFinalityIndexes []string `json:"fastFinalityIndexes"`
	DeadlineHeight      int64    `json:"deadlineHeight"`
	Block               int64    `json:"block"`
	ValidatorSet        []string `json:"validatorSet"`
	Timestamp           int64    `json:"timestamp"`
}

// UploadCost represents the cost estimate for uploading data
type UploadCost struct {
	Winc        string                 `json:"winc"`
	Bytes       int64                  `json:"bytes"`
	Adjustments map[string]interface{} `json:"adjustments,omitempty"`
}

// SignedDataItemUploadRequest represents a request to upload a pre-signed data item
type SignedDataItemUploadRequest struct {
	DataItemStreamFactory func() (io.ReadCloser, error) `json:"-"`
	DataItemSizeFactory   func() int64                  `json:"-"`
	Events                *UploadEvents                 `json:"-"`
	Context               context.Context               `json:"-"`
}
