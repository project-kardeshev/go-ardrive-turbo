package signers

import (
	"context"
	"fmt"

	goar "github.com/everFinance/goar"
	goarTypes "github.com/everFinance/goar/types"
	goether "github.com/everFinance/goether"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// EthereumSigner implements the Signer interface for Ethereum wallets
type EthereumSigner struct {
	wallet     string
	signer     goether.Signer
	itemSigner goar.ItemSigner
	Address    string
	PublicKey  string
}

// NewEthereumSigner creates a new Ethereum signer from a private key
func NewEthereumSigner(wallet string) (*EthereumSigner, error) {
	signer, signerErr := goether.NewSigner(wallet)
	if signerErr != nil {
		return nil, signerErr
	}

	itemSigner, itemSignerErr := goar.NewItemSigner(signer)
	if itemSignerErr != nil {
		return nil, itemSignerErr
	}

	return &EthereumSigner{
		wallet:     wallet,
		signer:     *signer,
		itemSigner: *itemSigner,
		Address:    signer.Address.String(),
		PublicKey:  signer.GetPublicKeyHex(),
	}, nil
}

// GetNativeAddress returns the Ethereum address of the wallet
func (e *EthereumSigner) GetNativeAddress() (string, error) {
	return e.Address, nil
}

// GetTokenType returns the Ethereum token type
func (e *EthereumSigner) GetTokenType() turboTypes.TokenType {
	return turboTypes.TokenTypeEthereum
}

// Sign signs the provided data using the Ethereum wallet
func (e *EthereumSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	// Use goether signer to sign the data
	signature, err := e.signer.SignMsg(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	return signature, nil
}

// SignDataItem signs a data item and returns the signed bundle item
func (e *EthereumSigner) SignDataItem(ctx context.Context, dataItem *DataItem) (goarTypes.BundleItem, error) {
	// Convert our tags to goar tags
	goarTags := make([]goarTypes.Tag, len(dataItem.Tags))
	for i, tag := range dataItem.Tags {
		goarTags[i] = goarTypes.Tag{
			Name:  tag.Name,
			Value: tag.Value,
		}
	}

	// Use ItemSigner to create and sign the data item properly
	bundleItem, err := e.itemSigner.CreateAndSignItem(
		dataItem.Data,
		dataItem.Target,
		dataItem.Anchor,
		goarTags,
	)
	if err != nil {
		return goarTypes.BundleItem{}, fmt.Errorf("failed to create and sign data item: %w", err)
	}

	// ItemBinary should now be properly populated by CreateAndSignItem
	if len(bundleItem.ItemBinary) == 0 {
		return goarTypes.BundleItem{}, fmt.Errorf("failed to generate signed data item binary")
	}

	return bundleItem, nil
}
