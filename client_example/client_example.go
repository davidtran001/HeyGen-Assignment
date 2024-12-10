package main

import (
	"context"
	"log"
	"time"

	"github.com/davidtran001/HeyGen-Assignment/hgclient"
)

func main() {
	client := hgclient.NewClient("http://localhost:8080", 1, 5, 30, 2.0)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := client.WaitForResult(ctx)
	if err != nil {
		log.Println("Failed to get result:", err.Error())
		return
	}

	log.Println("Final result is:", result)
}
