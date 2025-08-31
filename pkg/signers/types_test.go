package signers

import (
	"context"
	"errors"
	"testing"

	"github.com/everFinance/goar/types"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func TestCreateDataItem(t *testing.T) {
	data := []byte("test data")
	tags := []turboTypes.Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "test-app"},
	}
	target := "test-target"
	anchor := "test-anchor"

	dataItem := CreateDataItem(data, tags, target, anchor)

	if dataItem == nil {
		t.Error("Expected non-nil DataItem")
	}

	if string(dataItem.Data) != string(data) {
		t.Errorf("Expected data '%s', got '%s'", string(data), string(dataItem.Data))
	}

	if len(dataItem.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(dataItem.Tags))
	}

	for i, tag := range dataItem.Tags {
		if tag.Name != tags[i].Name {
			t.Errorf("Expected tag name '%s', got '%s'", tags[i].Name, tag.Name)
		}
		if tag.Value != tags[i].Value {
			t.Errorf("Expected tag value '%s', got '%s'", tags[i].Value, tag.Value)
		}
	}

	if dataItem.Target != target {
		t.Errorf("Expected target '%s', got '%s'", target, dataItem.Target)
	}

	if dataItem.Anchor != anchor {
		t.Errorf("Expected anchor '%s', got '%s'", anchor, dataItem.Anchor)
	}
}

func TestMockSigner(t *testing.T) {
	address := "test-address"
	tokenType := turboTypes.TokenTypeArweave
	mockSigner := NewMockSigner(address, tokenType)

	// Test GetNativeAddress
	addr, err := mockSigner.GetNativeAddress()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if addr != address {
		t.Errorf("Expected address '%s', got '%s'", address, addr)
	}

	// Test GetTokenType
	tt := mockSigner.GetTokenType()
	if tt != tokenType {
		t.Errorf("Expected token type '%s', got '%s'", tokenType, tt)
	}

	// Test Sign
	ctx := context.Background()
	data := []byte("test data")
	signature, err := mockSigner.Sign(ctx, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(signature) != "mock-signature" {
		t.Errorf("Expected signature 'mock-signature', got '%s'", string(signature))
	}

	// Test SignDataItem
	dataItem := &DataItem{
		Data: data,
		Tags: []turboTypes.Tag{{Name: "test", Value: "value"}},
	}
	bundleItem, err := mockSigner.SignDataItem(ctx, dataItem)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(bundleItem.ItemBinary) != "mock-signed-data-item" {
		t.Errorf("Expected ItemBinary 'mock-signed-data-item', got '%s'", string(bundleItem.ItemBinary))
	}
}

func TestMockSignerErrors(t *testing.T) {
	mockSigner := NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	ctx := context.Background()

	// Test Sign error
	expectedSignError := errors.New("sign error")
	mockSigner.SetSignError(expectedSignError)

	_, err := mockSigner.Sign(ctx, []byte("test"))
	if err != expectedSignError {
		t.Errorf("Expected sign error '%v', got '%v'", expectedSignError, err)
	}

	// Test SignDataItem error
	expectedSignDataItemError := errors.New("sign data item error")
	mockSigner.SetSignDataItemError(expectedSignDataItemError)

	dataItem := &DataItem{Data: []byte("test")}
	_, err = mockSigner.SignDataItem(ctx, dataItem)
	if err != expectedSignDataItemError {
		t.Errorf("Expected sign data item error '%v', got '%v'", expectedSignDataItemError, err)
	}
}

func TestMockSignerCustomResults(t *testing.T) {
	mockSigner := NewMockSigner("test-address", turboTypes.TokenTypeEthereum)
	ctx := context.Background()

	// Test custom sign result
	customSignature := []byte("custom-signature")
	mockSigner.SetSignResult(customSignature)

	signature, err := mockSigner.Sign(ctx, []byte("test"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(signature) != string(customSignature) {
		t.Errorf("Expected signature '%s', got '%s'", string(customSignature), string(signature))
	}

	// Test custom sign data item result
	customBundleItem := types.BundleItem{
		ItemBinary: []byte("custom-signed-item"),
	}
	mockSigner.SetSignDataItemResult(customBundleItem)

	dataItem := &DataItem{Data: []byte("test")}
	bundleItem, err := mockSigner.SignDataItem(ctx, dataItem)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(bundleItem.ItemBinary) != string(customBundleItem.ItemBinary) {
		t.Errorf("Expected ItemBinary '%s', got '%s'",
			string(customBundleItem.ItemBinary), string(bundleItem.ItemBinary))
	}
}

func TestDataItemValidation(t *testing.T) {
	// Test with nil data
	dataItem := CreateDataItem(nil, nil, "", "")
	if dataItem == nil {
		t.Error("Expected non-nil DataItem even with nil data")
	}

	if dataItem.Data != nil {
		t.Error("Expected nil data")
	}

	if dataItem.Tags != nil {
		t.Error("Expected nil tags")
	}

	// Test with empty tags slice
	emptyTags := []turboTypes.Tag{}
	dataItem = CreateDataItem([]byte("test"), emptyTags, "", "")
	if len(dataItem.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(dataItem.Tags))
	}
}
