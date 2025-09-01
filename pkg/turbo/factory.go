package turbo

import (
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
)

// TurboFactory provides factory methods for creating Turbo clients
type TurboFactory struct{}

// TurboConfig contains configuration options for creating Turbo clients
type TurboConfig struct {
	PaymentURL string // Payment service URL
	UploadURL  string // Upload service URL
}

// DefaultConfig returns the default production configuration
func DefaultConfig() *TurboConfig {
	return &TurboConfig{
		PaymentURL: "https://payment.ardrive.io",
		UploadURL:  "https://upload.ardrive.io",
	}
}

// DevConfig returns the development configuration
func DevConfig() *TurboConfig {
	return &TurboConfig{
		PaymentURL: "https://payment.ardrive.dev",
		UploadURL:  "https://upload.ardrive.dev",
	}
}

// Unauthenticated creates a new unauthenticated Turbo client
func (f *TurboFactory) Unauthenticated(config *TurboConfig) TurboUnauthenticatedClient {
	if config == nil {
		config = DefaultConfig()
	}

	return NewUnauthenticatedClient(config.PaymentURL, config.UploadURL)
}

// Authenticated creates a new authenticated Turbo client with the provided signer
func (f *TurboFactory) Authenticated(config *TurboConfig, signer signers.Signer) TurboAuthenticatedClient {
	if config == nil {
		config = DefaultConfig()
	}

	return NewAuthenticatedClient(config.PaymentURL, config.UploadURL, signer)
}

// Global factory instance
var Factory = &TurboFactory{}

// Convenience functions for easier usage

// Unauthenticated creates a new unauthenticated Turbo client using the global factory
func Unauthenticated(config *TurboConfig) TurboUnauthenticatedClient {
	return Factory.Unauthenticated(config)
}

// Authenticated creates a new authenticated Turbo client using the global factory
func Authenticated(config *TurboConfig, signer signers.Signer) TurboAuthenticatedClient {
	return Factory.Authenticated(config, signer)
}
