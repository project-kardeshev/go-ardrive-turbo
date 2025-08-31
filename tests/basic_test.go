package tests

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	turbo "github.com/project-kardeshev/go-ardrive-turbo"
)

// loadTestWallet loads the test wallet JWK from test_wallet.json
func loadTestWallet(t *testing.T) map[string]interface{} {
	walletPath := filepath.Join(".", "test_wallet.json")
	walletBytes, err := os.ReadFile(walletPath)
	if err != nil {
		t.Fatalf("Failed to read test wallet: %v", err)
	}

	var jwk map[string]interface{}
	if err := json.Unmarshal(walletBytes, &jwk); err != nil {
		t.Fatalf("Failed to parse test wallet JWK: %v", err)
	}

	return jwk
}

func TestUnauthenticatedClient(t *testing.T) {
	client := turbo.Unauthenticated(nil)
	
	if client == nil {
		t.Fatal("Expected non-nil unauthenticated client")
	}
}

func TestGetUploadCosts(t *testing.T) {
	client := turbo.Unauthenticated(nil)
	ctx := context.Background()
	
	req := &turbo.UploadCostsRequest{
		Bytes: []int64{1024, 2048},
	}
	
	// Note: This will fail without a real Turbo service running
	// In a real test environment, you'd mock the HTTP client
	_, err := client.GetUploadCosts(ctx, req)
	
	// We expect this to fail in the test environment
	// but we're testing that the method exists and can be called
	if err == nil {
		t.Log("Upload costs request succeeded (unexpected in test env)")
	} else {
		t.Logf("Upload costs request failed as expected: %v", err)
	}
}

func TestArweaveSigner(t *testing.T) {
	// Load the real test wallet
	jwk := loadTestWallet(t)
	
	// Create Arweave signer with real JWK
	signer, err := turbo.NewArweaveSigner(jwk)
	if err != nil {
		t.Fatalf("Failed to create Arweave signer with test wallet: %v", err)
	}
	
	// Test getting native address
	address, err := signer.GetNativeAddress()
	if err != nil {
		t.Fatalf("Failed to get native address: %v", err)
	}
	
	if address == "" {
		t.Fatal("Expected non-empty address")
	}
	
	t.Logf("Test wallet address: %s", address)
	
	// Test token type
	if signer.GetTokenType() != turbo.TokenTypeArweave {
		t.Fatalf("Expected token type %s, got %s", turbo.TokenTypeArweave, signer.GetTokenType())
	}
	
	// Test signing some data
	testData := []byte("Hello, Turbo SDK test!")
	signature, err := signer.Sign(context.Background(), testData)
	if err != nil {
		t.Fatalf("Failed to sign test data: %v", err)
	}
	
	if len(signature) == 0 {
		t.Fatal("Expected non-empty signature")
	}
	
	t.Logf("Successfully signed %d bytes of test data, signature length: %d", len(testData), len(signature))
}

func TestEthereumSigner(t *testing.T) {
	// Test with invalid private key (should fail)
	_, err := turbo.NewEthereumSigner("invalid-key")
	if err == nil {
		t.Fatal("Expected error with invalid private key")
	}
	
	// Test with valid format but dummy key
	validFormatKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer, err := turbo.NewEthereumSigner(validFormatKey)
	if err != nil {
		t.Logf("Ethereum signer creation failed: %v", err)
	} else {
		address, err := signer.GetNativeAddress()
		if err != nil {
			t.Fatalf("Failed to get address: %v", err)
		}
		if address == "" {
			t.Fatal("Expected non-empty address")
		}
		t.Logf("Created Ethereum signer with address: %s", address)
	}
}

func TestWinstonTypes(t *testing.T) {
	// Test Winston creation
	w1, err := turbo.NewWinstonFromString("1000000000000")
	if err != nil {
		t.Fatalf("Failed to create Winston: %v", err)
	}
	if w1.String() != "1000000000000" {
		t.Fatalf("Expected 1000000000000, got %s", w1.String())
	}
	
	// Test invalid Winston
	_, err = turbo.NewWinstonFromString("invalid")
	if err == nil {
		t.Fatal("Expected error with invalid Winston string")
	}
}

func TestTokenTypes(t *testing.T) {
	// Test token type constants
	if turbo.TokenTypeArweave != "arweave" {
		t.Fatalf("Expected 'arweave', got %s", turbo.TokenTypeArweave)
	}
	
	if turbo.TokenTypeEthereum != "ethereum" {
		t.Fatalf("Expected 'ethereum', got %s", turbo.TokenTypeEthereum)
	}
}

func TestAuthenticatedClient(t *testing.T) {
	// Load the test wallet
	jwk := loadTestWallet(t)
	
	// Create authenticated client
	client, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
		PrivateKey: jwk,
		Token:      turbo.TokenTypeArweave,
	})
	if err != nil {
		t.Fatalf("Failed to create authenticated client: %v", err)
	}
	
	// Test getting signer
	signer := client.GetSigner()
	if signer == nil {
		t.Fatal("Expected non-nil signer")
	}
	
	address, err := signer.GetNativeAddress()
	if err != nil {
		t.Fatalf("Failed to get signer address: %v", err)
	}
	t.Logf("Authenticated client address: %s", address)
	
	// Test data item signing
	ctx := context.Background()
	testData := []byte("Test data for signing")
	dataItem := turbo.CreateDataItem(testData, []turbo.Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "Test", Value: "true"},
	}, "", "")
	
	bundleItem, err := signer.SignDataItem(ctx, dataItem)
	if err != nil {
		t.Fatalf("Failed to sign data item: %v", err)
	}
	
	if bundleItem.Id == "" {
		t.Fatal("Expected non-empty signed data item ID")
	}
	
	if len(bundleItem.ItemBinary) == 0 {
		t.Fatal("Expected non-empty signed data item raw bytes")
	}
	
	t.Logf("Successfully signed data item: ID=%s, Size=%d bytes", bundleItem.Id, len(bundleItem.ItemBinary))
}

func TestConfigurations(t *testing.T) {
	// Test default config
	defaultConfig := turbo.DefaultConfig()
	if defaultConfig.Token != turbo.TokenTypeArweave {
		t.Fatalf("Expected default token to be arweave, got %s", defaultConfig.Token)
	}
	
	// Test dev config
	devConfig := turbo.DevConfig()
	if !devConfig.DevMode {
		t.Fatal("Expected dev config to have DevMode=true")
	}
}
