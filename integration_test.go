package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/davidtran001/HeyGen-Assignment/hgclient"
	"github.com/stretchr/testify/require"
)

func TestIntegrationServer(t *testing.T) {
	go func() {
		// Call main() function from main.go which starts the server
		main()
	}()

	// Create an hgclient instance pointing to the server
	client := hgclient.NewClient("http://localhost:8080", 1, 5, 30, 2.0)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	result, err := client.WaitForResult(ctx)
	require.NoError(t, err, "should not error while waiting for result")

	// Result should be either 'completed' or 'error' after the server passes the pending stage
	require.True(t, result == "completed" || result == "error", "result should be completed or error")
	log.Println("Integration test finished. Result:", result)
}
