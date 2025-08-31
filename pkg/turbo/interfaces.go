package turbo

import (
	"context"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// TurboUnauthenticatedClient provides access to Turbo's unauthenticated services
type TurboUnauthenticatedClient interface {
	// GetBalance returns the credit balance for a given address
	GetBalance(ctx context.Context, address string) (*types.Balance, error)

	// GetUploadCosts returns the estimated cost in Winston Credits for the provided file sizes
	GetUploadCosts(ctx context.Context, bytes []int64) ([]types.UploadCost, error)

	// UploadSignedDataItem uploads a pre-signed data item
	UploadSignedDataItem(ctx context.Context, req *types.SignedDataItemUploadRequest) (*types.UploadResult, error)
}

// TurboAuthenticatedClient provides access to both authenticated and unauthenticated Turbo services
type TurboAuthenticatedClient interface {
	TurboUnauthenticatedClient

	// GetBalanceForSigner returns the credit balance of the authenticated wallet
	GetBalanceForSigner(ctx context.Context) (*types.Balance, error)

	// Upload signs and uploads data to Turbo
	Upload(ctx context.Context, req *types.UploadRequest) (*types.UploadResult, error)

	// GetSigner returns the signer associated with this client
	GetSigner() signers.Signer
}
