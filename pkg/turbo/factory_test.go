package turbo

import (
	"testing"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	turboTypes "github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Error("Expected non-nil config")
	}

	expectedPaymentURL := "https://payment.ardrive.io"
	if config.PaymentURL != expectedPaymentURL {
		t.Errorf("Expected payment URL '%s', got '%s'", expectedPaymentURL, config.PaymentURL)
	}

	expectedUploadURL := "https://upload.ardrive.io"
	if config.UploadURL != expectedUploadURL {
		t.Errorf("Expected upload URL '%s', got '%s'", expectedUploadURL, config.UploadURL)
	}
}

func TestDevConfig(t *testing.T) {
	config := DevConfig()

	if config == nil {
		t.Error("Expected non-nil config")
	}

	expectedPaymentURL := "https://payment.ardrive.dev"
	if config.PaymentURL != expectedPaymentURL {
		t.Errorf("Expected payment URL '%s', got '%s'", expectedPaymentURL, config.PaymentURL)
	}

	expectedUploadURL := "https://upload.ardrive.dev"
	if config.UploadURL != expectedUploadURL {
		t.Errorf("Expected upload URL '%s', got '%s'", expectedUploadURL, config.UploadURL)
	}
}

func TestTurboFactoryUnauthenticated(t *testing.T) {
	factory := &TurboFactory{}

	// Test with default config
	client := factory.Unauthenticated(nil)
	if client == nil {
		t.Error("Expected non-nil unauthenticated client")
	}

	// Test with custom config
	config := &TurboConfig{
		PaymentURL: "https://custom-payment.test",
		UploadURL:  "https://custom-upload.test",
	}

	client = factory.Unauthenticated(config)
	if client == nil {
		t.Error("Expected non-nil unauthenticated client with custom config")
	}
}

func TestTurboFactoryAuthenticated(t *testing.T) {
	factory := &TurboFactory{}
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)

	// Test with default config
	client := factory.Authenticated(nil, mockSigner)
	if client == nil {
		t.Error("Expected non-nil authenticated client")
	}

	if client.GetSigner() != mockSigner {
		t.Error("Expected signer to match")
	}

	// Test with custom config
	config := &TurboConfig{
		PaymentURL: "https://custom-payment.test",
		UploadURL:  "https://custom-upload.test",
	}

	client = factory.Authenticated(config, mockSigner)
	if client == nil {
		t.Error("Expected non-nil authenticated client with custom config")
	}

	if client.GetSigner() != mockSigner {
		t.Error("Expected signer to match with custom config")
	}
}

func TestGlobalFactoryFunctions(t *testing.T) {
	// Test global Unauthenticated function
	client := Unauthenticated(nil)
	if client == nil {
		t.Error("Expected non-nil client from global Unauthenticated function")
	}

	// Test global Authenticated function
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeEthereum)
	authClient := Authenticated(nil, mockSigner)
	if authClient == nil {
		t.Error("Expected non-nil client from global Authenticated function")
	}

	if authClient.GetSigner() != mockSigner {
		t.Error("Expected signer to match from global Authenticated function")
	}
}

func TestGlobalFactoryInstance(t *testing.T) {
	if Factory == nil {
		t.Error("Expected non-nil global Factory instance")
	}

	// Test that global functions use the Factory instance
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)

	client1 := Factory.Unauthenticated(nil)
	client2 := Unauthenticated(nil)

	// Both should be non-nil (can't easily test they're the same type without reflection)
	if client1 == nil {
		t.Error("Expected non-nil client from Factory.Unauthenticated")
	}

	if client2 == nil {
		t.Error("Expected non-nil client from global Unauthenticated")
	}

	authClient1 := Factory.Authenticated(nil, mockSigner)
	authClient2 := Authenticated(nil, mockSigner)

	if authClient1 == nil {
		t.Error("Expected non-nil client from Factory.Authenticated")
	}

	if authClient2 == nil {
		t.Error("Expected non-nil client from global Authenticated")
	}
}

func TestTurboConfigCustomValues(t *testing.T) {
	config := &TurboConfig{
		PaymentURL: "https://my-payment.custom",
		UploadURL:  "https://my-upload.custom",
	}

	factory := &TurboFactory{}

	// Test unauthenticated client with custom config
	unauthClient := factory.Unauthenticated(config)
	if unauthClient == nil {
		t.Error("Expected non-nil unauthenticated client")
	}

	// Test authenticated client with custom config
	mockSigner := signers.NewMockSigner("test-address", turboTypes.TokenTypeArweave)
	authClient := factory.Authenticated(config, mockSigner)
	if authClient == nil {
		t.Error("Expected non-nil authenticated client")
	}

	if authClient.GetSigner() != mockSigner {
		t.Error("Expected signer to match")
	}
}

func TestConfigImmutability(t *testing.T) {
	// Test that default configs return new instances
	config1 := DefaultConfig()
	config2 := DefaultConfig()

	if config1 == config2 {
		t.Error("Expected different instances of default config")
	}

	// Test that dev configs return new instances
	devConfig1 := DevConfig()
	devConfig2 := DevConfig()

	if devConfig1 == devConfig2 {
		t.Error("Expected different instances of dev config")
	}

	// Modify one config and ensure the other is unchanged
	originalPaymentURL := config1.PaymentURL
	config1.PaymentURL = "https://modified.test"

	if config2.PaymentURL != originalPaymentURL {
		t.Error("Expected config2 to be unaffected by config1 modification")
	}
}
