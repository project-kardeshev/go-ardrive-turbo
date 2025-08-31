package signers

import (
	"context"

	"github.com/everFinance/goar/types"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// MockSigner implements the Signer interface for testing
type MockSigner struct {
	Address            string
	TokenType          turboTypes.TokenType
	SignError          error
	SignDataItemError  error
	SignResult         []byte
	SignDataItemResult types.BundleItem
}

// NewMockSigner creates a new mock signer
func NewMockSigner(address string, tokenType turboTypes.TokenType) *MockSigner {
	return &MockSigner{
		Address:    address,
		TokenType:  tokenType,
		SignResult: []byte("mock-signature"),
		SignDataItemResult: types.BundleItem{
			ItemBinary: []byte("mock-signed-data-item"),
		},
	}
}

// GetNativeAddress returns the mock address
func (m *MockSigner) GetNativeAddress() (string, error) {
	return m.Address, nil
}

// GetTokenType returns the mock token type
func (m *MockSigner) GetTokenType() turboTypes.TokenType {
	return m.TokenType
}

// Sign returns a mock signature or error
func (m *MockSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	if m.SignError != nil {
		return nil, m.SignError
	}
	return m.SignResult, nil
}

// SignDataItem returns a mock signed data item or error
func (m *MockSigner) SignDataItem(ctx context.Context, dataItem *DataItem) (types.BundleItem, error) {
	if m.SignDataItemError != nil {
		return types.BundleItem{}, m.SignDataItemError
	}
	return m.SignDataItemResult, nil
}

// SetSignError sets an error to be returned by Sign
func (m *MockSigner) SetSignError(err error) {
	m.SignError = err
}

// SetSignDataItemError sets an error to be returned by SignDataItem
func (m *MockSigner) SetSignDataItemError(err error) {
	m.SignDataItemError = err
}

// SetSignResult sets the result to be returned by Sign
func (m *MockSigner) SetSignResult(result []byte) {
	m.SignResult = result
}

// SetSignDataItemResult sets the result to be returned by SignDataItem
func (m *MockSigner) SetSignDataItemResult(result types.BundleItem) {
	m.SignDataItemResult = result
}
