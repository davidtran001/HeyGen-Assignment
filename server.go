package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Result string `json:"result"`
}

func main() {
	delayMs := 10000 // in ms

	// Capture the start time
	startTime := time.Now()

	r := chi.NewRouter()
	//r.Use(middleware.Logger)

	// GET: /status will send the status of the current request
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {

		elapsed := time.Since(startTime)
		// Request is still 'processing' so return pending
		if elapsed.Milliseconds() < int64(delayMs) {
			resp := Response{Result: "pending"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// After the delay, return either completed or error
		completionRate := 9
		if rand.Intn(10) < completionRate { // Request will complete at a percentage of: (completionRate / 10) * 100 [%]
			resp := Response{Result: "completed"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := Response{Result: "error"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}

		// Reset the timer when request is completed or error
		startTime = time.Now()
	})

	log.Printf("Starting server on :%s\n", "8080")
	http.ListenAndServe(":8080", r)
}
