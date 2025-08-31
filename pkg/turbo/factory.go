package turbo

import (
	"fmt"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

// TurboFactory provides factory methods for creating Turbo clients
type TurboFactory struct{}

// UnauthenticatedOptions contains options for creating an unauthenticated client
type UnauthenticatedOptions struct {
	Config *types.Config
}

// AuthenticatedOptions contains options for creating an authenticated client
type AuthenticatedOptions struct {
	Signer     signers.Signer
	PrivateKey interface{} // Can be JWK map, hex string, etc.
	Token      types.TokenType
	Config     *types.Config
}

// Unauthenticated creates a new unauthenticated Turbo client
func (f *TurboFactory) Unauthenticated(opts *UnauthenticatedOptions) TurboUnauthenticatedClient {
	config := types.DefaultConfig()
	if opts != nil && opts.Config != nil {
		config = opts.Config
	}

	return &unauthenticatedClient{
		config:     config,
		httpClient: newHTTPClient(config),
	}
}

// Authenticated creates a new authenticated Turbo client
func (f *TurboFactory) Authenticated(opts *AuthenticatedOptions) (TurboAuthenticatedClient, error) {
	if opts == nil {
		return nil, fmt.Errorf("authenticated options are required")
	}

	config := types.DefaultConfig()
	if opts.Config != nil {
		config = opts.Config
	}

	var signer signers.Signer
	var err error

	// If signer is provided, use it directly
	if opts.Signer != nil {
		signer = opts.Signer
	} else if opts.PrivateKey != nil {
		// Create signer based on token type and private key
		signer, err = createSignerFromPrivateKey(opts.PrivateKey, opts.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to create signer: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either signer or private key must be provided")
	}

	// Update config token type based on signer
	config.Token = signer.GetTokenType()

	return &authenticatedClient{
		unauthenticatedClient: &unauthenticatedClient{
			config:     config,
			httpClient: newHTTPClient(config),
		},
		signer: signer,
	}, nil
}

// createSignerFromPrivateKey creates a signer based on the token type and private key
func createSignerFromPrivateKey(privateKey interface{}, tokenType types.TokenType) (signers.Signer, error) {
	switch tokenType {
	case types.TokenTypeArweave:
		jwk, ok := privateKey.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("arweave private key must be a JWK map")
		}
		return signers.NewArweaveSigner(jwk)
		
	case types.TokenTypeEthereum:
		keyStr, ok := privateKey.(string)
		if !ok {
			return nil, fmt.Errorf("ethereum private key must be a hex string")
		}
		return signers.NewEthereumSigner(keyStr)
		
	default:
		return nil, fmt.Errorf("unsupported token type: %s", tokenType)
	}
}

// Global factory instance
var Factory = &TurboFactory{}

// Convenience functions
func Unauthenticated(opts *UnauthenticatedOptions) TurboUnauthenticatedClient {
	return Factory.Unauthenticated(opts)
}

func Authenticated(opts *AuthenticatedOptions) (TurboAuthenticatedClient, error) {
	return Factory.Authenticated(opts)
}
