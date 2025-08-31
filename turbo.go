// Package turbo provides a Go SDK for interacting with the Turbo Upload and Payment Services.
//
// This SDK allows you to:
// - Get wallet balance in Winston Credits
// - Estimate upload costs for data
// - Upload data with automatic signing
// - Upload pre-signed data items
//
// Example usage:
//
//	// Create an unauthenticated client
//	client := turbo.Unauthenticated(nil)
//	costs, err := client.GetUploadCosts(ctx, &types.UploadCostsRequest{
//		Bytes: []int64{1024, 2048},
//	})
//
//	// Create an authenticated client with Arweave JWK
//	authClient, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
//		PrivateKey: jwkMap,
//		Token:      types.TokenTypeArweave,
//	})
//
//	// Get balance
//	balance, err := authClient.GetBalance(ctx)
//
//	// Upload data
//	result, err := authClient.Upload(ctx, &types.UploadRequest{
//		Data: []byte("Hello, Turbo!"),
//		Tags: []types.Tag{{Name: "Content-Type", Value: "text/plain"}},
//	})
package turbo

import (
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/turbo"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// Re-export main types and interfaces for convenience
type (
	// Client interfaces
	TurboUnauthenticatedClient = turbo.TurboUnauthenticatedClient
	TurboAuthenticatedClient   = turbo.TurboAuthenticatedClient
	
	// Configuration and options
	Config                 = types.Config
	UnauthenticatedOptions = turbo.UnauthenticatedOptions
	AuthenticatedOptions   = turbo.AuthenticatedOptions
	
	// Core types
	Winston              = types.Winston
	Balance              = types.Balance
	UploadCost           = types.UploadCost
	UploadCostsRequest   = types.UploadCostsRequest
	UploadCostsResponse  = types.UploadCostsResponse
	UploadResult         = types.UploadResult
	UploadRequest        = types.UploadRequest
	UploadEvents         = types.UploadEvents
	ProgressEvent        = types.ProgressEvent
	ErrorEvent           = types.ErrorEvent
	Tag                  = types.Tag
	TokenType            = types.TokenType
	
	// Signer types
	Signer           = signers.Signer
	ArweaveSigner    = signers.ArweaveSigner
	EthereumSigner   = signers.EthereumSigner
	DataItem         = signers.DataItem
)

// Re-export constants
const (
	TokenTypeArweave  = types.TokenTypeArweave
	TokenTypeEthereum = types.TokenTypeEthereum
	TokenTypeSolana   = types.TokenTypeSolana
	TokenTypePolygon  = types.TokenTypePolygon
	TokenTypeKyve     = types.TokenTypeKyve
	TokenTypeBaseEth  = types.TokenTypeBaseEth
	TokenTypeArio     = types.TokenTypeArio
)

// Re-export factory functions
var (
	Unauthenticated = turbo.Unauthenticated
	Authenticated   = turbo.Authenticated
	Factory         = turbo.Factory
)

// Re-export configuration functions
var (
	DefaultConfig = types.DefaultConfig
	DevConfig     = types.DevConfig
)

// Re-export signer constructors
var (
	NewArweaveSigner           = signers.NewArweaveSigner
	NewArweaveSignerFromKeyfile = signers.NewArweaveSignerFromKeyfile
	NewEthereumSigner          = signers.NewEthereumSigner
	CreateDataItem             = signers.CreateDataItem
	CreateDataItemFromReader   = signers.CreateDataItemFromReader
)

// Re-export utility functions
var (
	NewWinston           = types.NewWinston
	NewWinstonFromString = types.NewWinstonFromString
)
