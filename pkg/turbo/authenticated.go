package turbo

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// authenticatedClient implements TurboAuthenticatedClient
type authenticatedClient struct {
	*unauthenticatedClient
	signer signers.Signer
}

// GetBalance returns the credit balance of the authenticated wallet
func (a *authenticatedClient) GetBalance(ctx context.Context) (*types.Balance, error) {
	// Get the wallet address
	address, err := a.signer.GetNativeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet address: %w", err)
	}

	// Make the request (no auth headers needed for Phase 1)
	url := fmt.Sprintf("%s/v1/balance?address=%s", a.httpClient.GetPaymentURL(), address)
	resp, err := a.httpClient.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Parse the response
	var balance types.Balance
	if err := ParseJSON(resp, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// Upload signs and uploads data to Turbo
func (a *authenticatedClient) Upload(ctx context.Context, req *types.UploadRequest) (*types.UploadResult, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request is required")
	}

	// Determine data source
	var data []byte
	var err error

	if req.Data != nil {
		data = req.Data
	} else if req.DataReader != nil {
		data, err = io.ReadAll(req.DataReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either Data or DataReader must be provided")
	}

	// Create upload context
	uploadCtx := ctx
	if req.Context != nil {
		uploadCtx = req.Context
	}

	// Notify signing start
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     int64(len(data)),
			ProcessedBytes: 0,
			Step:           "signing",
		})
	}

	// Create data item
	dataItem := signers.CreateDataItem(data, req.Tags, req.Target, req.Anchor)

	// Sign the data item
	bundleItem, err := a.signer.SignDataItem(uploadCtx, dataItem)
	if err != nil {
		if req.Events != nil && req.Events.OnSigningError != nil {
			req.Events.OnSigningError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "signing"})
		}
		return nil, fmt.Errorf("failed to sign data item: %w", err)
	}

	// Notify signing success
	if req.Events != nil && req.Events.OnSigningSuccess != nil {
		req.Events.OnSigningSuccess()
	}
	if req.Events != nil && req.Events.OnProgress != nil {
		req.Events.OnProgress(types.ProgressEvent{
			TotalBytes:     int64(len(data)),
			ProcessedBytes: int64(len(data)),
			Step:           "signing",
		})
	}

	// Create upload request for signed data item
	uploadReq := &types.SignedDataItemUploadRequest{
		DataItemStreamFactory: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bundleItem.ItemBinary)), nil
		},
		DataItemSizeFactory: func() int64 {
			return int64(len(bundleItem.ItemBinary))
		},
		Events:  req.Events,
		Context: uploadCtx,
	}

	// Upload the signed data item
	return a.UploadSignedDataItem(uploadCtx, uploadReq)
}

// GetSigner returns the signer associated with this client
func (a *authenticatedClient) GetSigner() signers.Signer {
	return a.signer
}
