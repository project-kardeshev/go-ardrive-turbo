package turbo

import (
	"context"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// TurboUnauthenticatedClient represents the interface for unauthenticated Turbo operations
type TurboUnauthenticatedClient interface {
	// GetUploadCosts returns the estimated cost in Winston Credits for the provided file sizes
	GetUploadCosts(ctx context.Context, req *types.UploadCostsRequest) (*types.UploadCostsResponse, error)
	
	// UploadSignedDataItem uploads a signed data item to Turbo
	UploadSignedDataItem(ctx context.Context, req *types.SignedDataItemUploadRequest) (*types.UploadResult, error)
}

// TurboAuthenticatedClient represents the interface for authenticated Turbo operations
type TurboAuthenticatedClient interface {
	TurboUnauthenticatedClient
	
	// GetBalance returns the credit balance of the authenticated wallet
	GetBalance(ctx context.Context) (*types.Balance, error)
	
	// Upload signs and uploads data to Turbo
	Upload(ctx context.Context, req *types.UploadRequest) (*types.UploadResult, error)
	
	// GetSigner returns the signer associated with this client
	GetSigner() signers.Signer
}
