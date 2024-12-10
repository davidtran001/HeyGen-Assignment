package hgclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockResponse struct {
    Result string `json:"result"`
}

// Test the scenario where the server returns "completed" immediately
func TestWaitForResult_Completed(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        resp := mockResponse{Result: "completed"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }))
    defer ts.Close()

	// Create client to expect 'completed' result
    client := NewClient(ts.URL, 1, 5, 10, 2.0)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := client.WaitForResult(ctx)
    require.NoError(t, err, "should not return an error")
    require.Equal(t, "completed", result, "result should be completed immediately")
}

// Test the scenario where the server immediately returns "error"
func TestWaitForResult_Error(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        resp := mockResponse{Result: "error"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }))
    defer ts.Close()

	// Create client to expect 'error' result
    client := NewClient(ts.URL, 1, 5, 10, 2.0)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := client.WaitForResult(ctx)
    require.NoError(t, err)
    require.Equal(t, "error", result, "result should be error immediately")
}

// Test the scenario where the server returns "pending" a few times then "completed"
func TestWaitForResult_PendingThenCompleted(t *testing.T) {
    attemptCount := 0

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // For the first 2 requests, return pending
        // Afterwards, return completed
        if attemptCount < 2 {
            resp := mockResponse{Result: "pending"}
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(resp)
            attemptCount++
            return
        }

        resp := mockResponse{Result: "completed"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }))
    defer ts.Close()

	// Create client to expect 'pending' and then 'completed' result
    client := NewClient(ts.URL, 1, 5, 30, 2.0)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, err := client.WaitForResult(ctx)
    require.NoError(t, err)
    require.Equal(t, "completed", result, "after pending attempts, result should eventually be completed")
}

// Test that the client times out if pending never resolves
func TestWaitForResult_Timeout(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        resp := mockResponse{Result: "pending"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }))
    defer ts.Close()

    // The total waiting time is small in order to trigger a timeout quickly
    client := NewClient(ts.URL, 1, 1, 3, 2.0)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := client.WaitForResult(ctx)
    require.Error(t, err, "should return a timeout error")
    require.Contains(t, err.Error(), "timed out waiting for result")
    require.Empty(t, result, "no result expected on timeout")
}
