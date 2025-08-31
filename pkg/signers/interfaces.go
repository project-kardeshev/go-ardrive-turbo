package signers

import (
	"context"
	"io"

	goarTypes "github.com/everFinance/goar/types"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// Signer represents a generic signer interface
type Signer interface {
	// GetNativeAddress returns the native address of the signer
	GetNativeAddress() (string, error)
	
	// GetTokenType returns the token type this signer supports
	GetTokenType() types.TokenType
	
	// Sign signs the provided data and returns the signature
	Sign(ctx context.Context, data []byte) ([]byte, error)
	
	// SignDataItem signs a data item and returns the signed bundle item
	SignDataItem(ctx context.Context, dataItem *DataItem) (goarTypes.BundleItem, error)
}

// DataItem represents an unsigned Arweave data item
type DataItem struct {
	Data   []byte
	Tags   []types.Tag
	Target string
	Anchor string
}

// CreateDataItem creates a new data item from the provided parameters
func CreateDataItem(data []byte, tags []types.Tag, target, anchor string) *DataItem {
	return &DataItem{
		Data:   data,
		Tags:   tags,
		Target: target,
		Anchor: anchor,
	}
}

// CreateDataItemFromReader creates a new data item from a reader
func CreateDataItemFromReader(reader io.Reader, tags []types.Tag, target, anchor string) (*DataItem, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	
	return &DataItem{
		Data:   data,
		Tags:   tags,
		Target: target,
		Anchor: anchor,
	}, nil
}

