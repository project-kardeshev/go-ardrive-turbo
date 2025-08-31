package signers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// ArweaveSigner implements the Signer interface for Arweave wallets
type ArweaveSigner struct {
	signer     *goar.Signer
	itemSigner *goar.ItemSigner
}

// NewArweaveSigner creates a new Arweave signer from a JWK
func NewArweaveSigner(jwk map[string]interface{}) (*ArweaveSigner, error) {
	// Convert JWK map to JSON bytes
	jwkBytes, err := json.Marshal(jwk)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JWK: %w", err)
	}

	signer, err := goar.NewSigner(jwkBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer from JWK: %w", err)
	}

	itemSigner, err := goar.NewItemSigner(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create item signer: %w", err)
	}

	return &ArweaveSigner{
		signer:     signer,
		itemSigner: itemSigner,
	}, nil
}

// NewArweaveSignerFromKeyfile creates a new Arweave signer from a keyfile path
func NewArweaveSignerFromKeyfile(keyfile string) (*ArweaveSigner, error) {
	signer, err := goar.NewSignerFromPath(keyfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer from keyfile: %w", err)
	}

	itemSigner, err := goar.NewItemSigner(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create item signer: %w", err)
	}

	return &ArweaveSigner{
		signer:     signer,
		itemSigner: itemSigner,
	}, nil
}

// GetNativeAddress returns the Arweave address of the wallet
func (a *ArweaveSigner) GetNativeAddress() (string, error) {
	return a.signer.Address, nil
}

// GetTokenType returns the Arweave token type
func (a *ArweaveSigner) GetTokenType() turboTypes.TokenType {
	return turboTypes.TokenTypeArweave
}

// Sign signs the provided data using the Arweave wallet
func (a *ArweaveSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	// Sign the data using the goar signer
	signature, err := a.signer.SignMsg(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	return signature, nil
}

// SignDataItem signs a data item and returns the signed bundle item
func (a *ArweaveSigner) SignDataItem(ctx context.Context, dataItem *DataItem) (types.BundleItem, error) {
	// Convert our tags to goar tags
	goarTags := make([]types.Tag, len(dataItem.Tags))
	for i, tag := range dataItem.Tags {
		goarTags[i] = types.Tag{
			Name:  tag.Name,
			Value: tag.Value,
		}
	}

	// Use ItemSigner to create and sign the data item properly
	bundleItem, err := a.itemSigner.CreateAndSignItem(
		dataItem.Data,
		dataItem.Target,
		dataItem.Anchor,
		goarTags,
	)
	if err != nil {
		return types.BundleItem{}, fmt.Errorf("failed to create and sign data item: %w", err)
	}

	// ItemBinary should now be properly populated by CreateAndSignItem
	if len(bundleItem.ItemBinary) == 0 {
		return types.BundleItem{}, fmt.Errorf("failed to generate signed data item binary")
	}

	return bundleItem, nil
}
