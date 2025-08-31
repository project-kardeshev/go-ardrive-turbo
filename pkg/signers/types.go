package signers

import (
	"context"

	"github.com/everFinance/goar/types"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// Signer interface for different wallet types
type Signer interface {
	GetNativeAddress() (string, error)
	GetTokenType() turboTypes.TokenType
	Sign(ctx context.Context, data []byte) ([]byte, error)
	SignDataItem(ctx context.Context, dataItem *DataItem) (types.BundleItem, error)
}

// DataItem represents a data item to be signed and uploaded
type DataItem struct {
	Data   []byte           `json:"data"`
	Tags   []turboTypes.Tag `json:"tags"`
	Target string           `json:"target"`
	Anchor string           `json:"anchor"`
}

// CreateDataItem creates a new DataItem with the provided parameters
func CreateDataItem(data []byte, tags []turboTypes.Tag, target, anchor string) *DataItem {
	return &DataItem{
		Data:   data,
		Tags:   tags,
		Target: target,
		Anchor: anchor,
	}
}
