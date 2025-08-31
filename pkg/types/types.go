package types

import (
	"context"
	"io"
	"math/big"
)

// Winston represents Winston Credits (the smallest unit of AR)
type Winston struct {
	*big.Int
}

// NewWinston creates a new Winston from a big.Int
func NewWinston(value *big.Int) Winston {
	return Winston{Int: new(big.Int).Set(value)}
}

// NewWinstonFromString creates a Winston from a string representation
func NewWinstonFromString(value string) (Winston, error) {
	w := &big.Int{}
	_, ok := w.SetString(value, 10)
	if !ok {
		return Winston{}, ErrInvalidWinstonValue
	}
	return Winston{Int: w}, nil
}

// String returns the string representation of Winston
func (w Winston) String() string {
	if w.Int == nil {
		return "0"
	}
	return w.Int.String()
}

// Balance represents a wallet's credit balance
type Balance struct {
	Winc Winston `json:"winc"`
}

// UploadCost represents the cost estimation for an upload
type UploadCost struct {
	Winc        Winston                `json:"winc"`
	Adjustments map[string]interface{} `json:"adjustments,omitempty"`
}

// UploadCostsRequest represents a request for upload cost estimation
type UploadCostsRequest struct {
	Bytes []int64 `json:"bytes"`
}

// UploadCostsResponse represents the response from upload cost estimation
type UploadCostsResponse []UploadCost

// UploadResult represents the result of a successful upload
type UploadResult struct {
	ID                    string   `json:"id"`
	Owner                 string   `json:"owner"`
	DataCaches            []string `json:"dataCaches,omitempty"`
	FastFinalityIndexes   []string `json:"fastFinalityIndexes,omitempty"`
	DeadlineHeight        *int64   `json:"deadlineHeight,omitempty"`
	Block                 *int64   `json:"block,omitempty"`
	ValidatorSignatures   []string `json:"validatorSignatures,omitempty"`
	Verify                *string  `json:"verify,omitempty"`
}

// UploadRequest represents a request to upload data
type UploadRequest struct {
	Data         []byte                 `json:"-"`
	DataReader   io.Reader              `json:"-"`
	Tags         []Tag                  `json:"tags,omitempty"`
	Target       string                 `json:"target,omitempty"`
	Anchor       string                 `json:"anchor,omitempty"`
	Events       *UploadEvents          `json:"-"`
	Context      context.Context        `json:"-"`
}

// Tag represents an Arweave tag
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UploadEvents contains callback functions for upload progress tracking
type UploadEvents struct {
	OnProgress        func(ProgressEvent)
	OnError           func(ErrorEvent)
	OnSuccess         func()
	OnSigningProgress func(ProgressEvent)
	OnSigningError    func(error)
	OnSigningSuccess  func()
	OnUploadProgress  func(ProgressEvent)
	OnUploadError     func(error)
	OnUploadSuccess   func()
}

// ProgressEvent represents progress information
type ProgressEvent struct {
	TotalBytes     int64  `json:"totalBytes"`
	ProcessedBytes int64  `json:"processedBytes"`
	Step           string `json:"step,omitempty"`
}

// ErrorEvent represents error information with context
type ErrorEvent struct {
	Error error  `json:"error"`
	Step  string `json:"step,omitempty"`
}

// SignedDataItemUploadRequest represents a request to upload a pre-signed data item
type SignedDataItemUploadRequest struct {
	DataItemStreamFactory func() (io.ReadCloser, error)
	DataItemSizeFactory   func() int64
	Events                *UploadEvents
	Context               context.Context
}

// TokenType represents the supported token types
type TokenType string

const (
	TokenTypeArweave  TokenType = "arweave"
	TokenTypeEthereum TokenType = "ethereum"
	TokenTypeSolana   TokenType = "solana"
	TokenTypePolygon  TokenType = "pol"
	TokenTypeKyve     TokenType = "kyve"
	TokenTypeBaseEth  TokenType = "base-eth"
	TokenTypeArio     TokenType = "ario"
)

// Config represents the configuration for the Turbo client
type Config struct {
	GatewayURL    string
	UploadURL     string
	PaymentURL    string
	Token         TokenType
	DevMode       bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		GatewayURL:  "https://arweave.net",
		UploadURL:   "https://upload.ardrive.io",
		PaymentURL:  "https://payment.ardrive.io",
		Token:       TokenTypeArweave,
		DevMode:     false,
	}
}

// DevConfig returns the development configuration
func DevConfig() *Config {
	return &Config{
		GatewayURL:  "https://arweave.net",
		UploadURL:   "https://upload.ardrive.dev",
		PaymentURL:  "https://payment.ardrive.dev",
		Token:       TokenTypeArweave,
		DevMode:     true,
	}
}
