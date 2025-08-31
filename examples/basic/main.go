package main

import (
	"context"
	"fmt"
	"log"

	turbo "github.com/project-kardeshev/go-ardrive-turbo"
)

func main() {
	ctx := context.Background()

	// Example 1: Get upload costs (unauthenticated)
	fmt.Println("=== Getting Upload Costs ===")
	client := turbo.Unauthenticated(nil)
	
	costs, err := client.GetUploadCosts(ctx, &turbo.UploadCostsRequest{
		Bytes: []int64{1024, 2048, 4096}, // 1KB, 2KB, 4KB
	})
	if err != nil {
		log.Printf("Failed to get upload costs: %v", err)
	} else {
		for i, cost := range *costs {
			fmt.Printf("Cost for %d bytes: %s winc\n", []int64{1024, 2048, 4096}[i], cost.Winc.String())
		}
	}

	// Example 2: Authenticated operations (using test wallet)
	// Uncomment to test with the provided test wallet
	/*
	import (
		"encoding/json"
		"os"
	)
	
	fmt.Println("\n=== Authenticated Operations ===")
	
	// Load test wallet (you can replace this with your own JWK)
	walletBytes, err := os.ReadFile("tests/test_wallet.json")
	if err != nil {
		log.Printf("Failed to read test wallet: %v", err)
		log.Printf("Make sure to run from the project root directory")
		return
	}
	
	var jwk map[string]interface{}
	if err := json.Unmarshal(walletBytes, &jwk); err != nil {
		log.Printf("Failed to parse test wallet: %v", err)
		return
	}
	
	authClient, err := turbo.Authenticated(&turbo.AuthenticatedOptions{
		PrivateKey: jwk,
		Token:      turbo.TokenTypeArweave,
	})
	if err != nil {
		log.Fatalf("Failed to create authenticated client: %v", err)
	}

	// Get balance
	balance, err := authClient.GetBalance(ctx)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
	} else {
		fmt.Printf("Current balance: %s winc\n", balance.Winc.String())
	}

	// Upload data
	uploadResult, err := authClient.Upload(ctx, &turbo.UploadRequest{
		Data: []byte("Hello from Go Turbo SDK!"),
		Tags: []turbo.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "App-Name", Value: "go-turbo-sdk-example"},
		},
		Events: &turbo.UploadEvents{
			OnProgress: func(event turbo.ProgressEvent) {
				fmt.Printf("Progress: %d/%d bytes (%s)\n", 
					event.ProcessedBytes, event.TotalBytes, event.Step)
			},
			OnError: func(event turbo.ErrorEvent) {
				fmt.Printf("Error during %s: %v\n", event.Step, event.Error)
			},
			OnSuccess: func() {
				fmt.Println("Upload completed successfully!")
			},
		},
	})
	if err != nil {
		log.Printf("Failed to upload: %v", err)
	} else {
		fmt.Printf("Upload successful! ID: %s\n", uploadResult.ID)
	}
	*/

	fmt.Println("\n=== Example completed ===")
	fmt.Println("To test authenticated operations, uncomment the code above and provide real wallet credentials.")
}
