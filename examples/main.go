package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/project-kardeshev/go-ardrive-turbo/pkg/signers"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/turbo"
	"github.com/project-kardeshev/go-ardrive-turbo/pkg/types"
)

func main() {
	ctx := context.Background()

	// Example 1: Unauthenticated client usage
	fmt.Println("=== Unauthenticated Client Examples ===")

	// Create unauthenticated client
	unauthClient := turbo.Unauthenticated(turbo.DevConfig())

	// Get balance for a specific address
	exampleAddress := "example-arweave-address"
	balance, err := unauthClient.GetBalance(ctx, exampleAddress)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
	} else {
		fmt.Printf("Balance for %s: %s winc\n", exampleAddress, balance.WinC)
	}

	// Get upload costs for different file sizes
	fileSizes := []int64{1024, 1024 * 1024, 10 * 1024 * 1024} // 1KB, 1MB, 10MB
	costs, err := unauthClient.GetUploadCosts(ctx, fileSizes)
	if err != nil {
		log.Printf("Failed to get upload costs: %v", err)
	} else {
		for i, cost := range costs {
			fmt.Printf("Cost for %d bytes: %s winc\n", fileSizes[i], cost.Winc)
		}
	}

	// Upload a pre-signed data item (example - would need actual signed data)
	signedDataItemUpload := &types.SignedDataItemUploadRequest{
		DataItemStreamFactory: func() (io.ReadCloser, error) {
			// This would be your actual signed data item binary
			exampleData := []byte("example-signed-data-item-binary")
			return io.NopCloser(bytes.NewReader(exampleData)), nil
		},
		DataItemSizeFactory: func() int64 {
			return int64(len("example-signed-data-item-binary"))
		},
		Events: &types.UploadEvents{
			OnProgress: func(event types.ProgressEvent) {
				fmt.Printf("Upload progress: %d/%d bytes (%s)\n",
					event.ProcessedBytes, event.TotalBytes, event.Step)
			},
			OnUploadSuccess: func(result *types.UploadResult) {
				fmt.Printf("Upload successful! ID: %s\n", result.ID)
			},
			OnUploadError: func(err error) {
				log.Printf("Upload error: %v", err)
			},
		},
	}

	result, err := unauthClient.UploadSignedDataItem(ctx, signedDataItemUpload)
	if err != nil {
		log.Printf("Failed to upload signed data item: %v", err)
	} else {
		fmt.Printf("Upload result: %+v\n", result)
	}

	fmt.Println("\n=== Authenticated Client Examples ===")

	// Example 2: Authenticated client usage (would need actual wallet/signer)
	// Note: This is just demonstrating the API, would need actual wallet keys

	// Example creating an Arweave signer (commented out since we don't have actual keys)
	/*
		jwk := map[string]interface{}{
			// Your JWK data here
		}
		arweaveSigner, err := signers.NewArweaveSigner(jwk)
		if err != nil {
			log.Fatalf("Failed to create Arweave signer: %v", err)
		}

		// Create authenticated client
		authClient := turbo.Authenticated(turbo.DevConfig(), arweaveSigner)

		// Get balance for the authenticated wallet
		balance, err := authClient.GetBalanceForSigner(ctx)
		if err != nil {
			log.Printf("Failed to get balance: %v", err)
		} else {
			fmt.Printf("Authenticated wallet balance: %s winc\n", balance.WinC)
		}

		// Upload data (signs and uploads)
		uploadRequest := &types.UploadRequest{
			Data: []byte("Hello, Turbo from Go!"),
			Tags: []types.Tag{
				{Name: "Content-Type", Value: "text/plain"},
				{Name: "App-Name", Value: "go-ardrive-turbo-example"},
			},
			Events: &types.UploadEvents{
				OnProgress: func(event types.ProgressEvent) {
					fmt.Printf("Upload progress: %d/%d bytes (%s)\n",
						event.ProcessedBytes, event.TotalBytes, event.Step)
				},
				OnSigningSuccess: func() {
					fmt.Println("Data signing successful!")
				},
				OnUploadSuccess: func(result *types.UploadResult) {
					fmt.Printf("Upload successful! ID: %s\n", result.ID)
				},
			},
		}

		result, err := authClient.Upload(ctx, uploadRequest)
		if err != nil {
			log.Printf("Failed to upload: %v", err)
		} else {
			fmt.Printf("Upload result: %+v\n", result)
		}
	*/

	fmt.Println("Example completed! (Authenticated examples commented out - need actual wallet keys)")
}

// demonstrateSignerCreation shows how to create different types of signers
func demonstrateSignerCreation() {
	fmt.Println("\n=== Signer Creation Examples ===")

	// Arweave signer from JWK
	jwk := map[string]interface{}{
		// Your JWK JSON data would go here
	}
	_, err := signers.NewArweaveSigner(jwk)
	if err != nil {
		fmt.Printf("Would create Arweave signer from JWK: %v\n", err)
	}

	// Arweave signer from keyfile
	keyfilePath := "/path/to/arweave-keyfile.json"
	_, err = signers.NewArweaveSignerFromKeyfile(keyfilePath)
	if err != nil {
		fmt.Printf("Would create Arweave signer from keyfile: %v\n", err)
	}

	// Ethereum signer
	ethPrivateKey := "your-ethereum-private-key-hex"
	_, err = signers.NewEthereumSigner(ethPrivateKey)
	if err != nil {
		fmt.Printf("Would create Ethereum signer: %v\n", err)
	}

	fmt.Println("Signer examples completed!")
}
