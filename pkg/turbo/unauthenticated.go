package turbo

import (
	"context"
	"fmt"
	"io"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// unauthenticatedClient implements TurboUnauthenticatedClient
type unauthenticatedClient struct {
	config     *types.Config
	httpClient *HTTPClient
}

// GetUploadCosts returns the estimated cost in Winston Credits for the provided file sizes
func (u *unauthenticatedClient) GetUploadCosts(ctx context.Context, req *types.UploadCostsRequest) (*types.UploadCostsResponse, error) {
	url := fmt.Sprintf("%s/v1/price/bytes", u.httpClient.GetPaymentURL())
	
	resp, err := u.httpClient.Post(ctx, url, req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload costs: %w", err)
	}

	var result types.UploadCostsResponse
	if err := ParseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UploadSignedDataItem uploads a signed data item to Turbo
func (u *unauthenticatedClient) UploadSignedDataItem(ctx context.Context, req *types.SignedDataItemUploadRequest) (*types.UploadResult, error) {
	if req.DataItemStreamFactory == nil {
		return nil, fmt.Errorf("data item stream factory is required")
	}

	if req.DataItemSizeFactory == nil {
		return nil, fmt.Errorf("data item size factory is required")
	}

	// Get the data item stream
	dataStream, err := req.DataItemStreamFactory()
	if err != nil {
		return nil, fmt.Errorf("failed to create data item stream: %w", err)
	}
	defer dataStream.Close()

	// Get the data item size
	dataSize := req.DataItemSizeFactory()

	// Create upload context
	uploadCtx := ctx
	if req.Context != nil {
		uploadCtx = req.Context
	}

	// Track progress if events are provided
	var progressReader io.Reader = dataStream
	if req.Events != nil && req.Events.OnUploadProgress != nil {
		progressReader = &progressReaderImpl{
			reader:      dataStream,
			totalBytes:  dataSize,
			onProgress:  req.Events.OnUploadProgress,
		}
	}

	// Prepare headers
	headers := map[string]string{
		"Content-Length": fmt.Sprintf("%d", dataSize),
	}

	// Upload the signed data item
	url := fmt.Sprintf("%s/v1/tx", u.httpClient.GetUploadURL())
	resp, err := u.httpClient.PostStream(uploadCtx, url, progressReader, "application/octet-stream", headers)
	if err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "upload"})
		}
		return nil, fmt.Errorf("failed to upload signed data item: %w", err)
	}

	// Parse the upload result
	var result types.UploadResult
	if err := ParseJSON(resp, &result); err != nil {
		if req.Events != nil && req.Events.OnUploadError != nil {
			req.Events.OnUploadError(err)
		}
		if req.Events != nil && req.Events.OnError != nil {
			req.Events.OnError(types.ErrorEvent{Error: err, Step: "upload"})
		}
		return nil, err
	}

	// Notify success
	if req.Events != nil && req.Events.OnUploadSuccess != nil {
		req.Events.OnUploadSuccess()
	}
	if req.Events != nil && req.Events.OnSuccess != nil {
		req.Events.OnSuccess()
	}

	return &result, nil
}

// progressReaderImpl wraps an io.Reader to track upload progress
type progressReaderImpl struct {
	reader         io.Reader
	totalBytes     int64
	processedBytes int64
	onProgress     func(types.ProgressEvent)
}

func (p *progressReaderImpl) Read(buf []byte) (int, error) {
	n, err := p.reader.Read(buf)
	if n > 0 {
		p.processedBytes += int64(n)
		if p.onProgress != nil {
			p.onProgress(types.ProgressEvent{
				TotalBytes:     p.totalBytes,
				ProcessedBytes: p.processedBytes,
				Step:           "upload",
			})
		}
	}
	return n, err
}
