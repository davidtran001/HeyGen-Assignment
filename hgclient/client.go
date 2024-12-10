package hgclient

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL       string
	HTTPClient    *http.Client
	InitialWait   time.Duration
	MaxWait       time.Duration
	MaxTotalTime  time.Duration
	BackoffFactor float64
}

type StatusResponse struct {
	Result string `json:"result"`
}

func NewClient(baseURL string, initialWait time.Duration, maxWait time.Duration, maxTotalTime time.Duration, backoffFactor float64) *Client {
	return &Client{
		BaseURL:       baseURL,
		HTTPClient:    &http.Client{Timeout: 5 * time.Second},
		InitialWait:   initialWait * time.Second,  // initial polling interval
		MaxWait:       maxWait * time.Second,      // max interval between polls
		MaxTotalTime:  maxTotalTime * time.Second, // total time before timeout
		BackoffFactor: backoffFactor,              // exponential backoff multiplier
	}
}

// GetStatus makes a simple HTTP GET request to the /status endpoint
func (c *Client) GetStatus(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/status", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status code from server")
	}

	var sr StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", err
	}
	return sr.Result, nil
}

// WaitForResult polls the server until the result is either completed, error, timeout
func (c *Client) WaitForResult(ctx context.Context) (string, error) {
	start := time.Now()
	interval := c.InitialWait

	// Poll the status request using exponential backoff
	for {
		if time.Since(start) > c.MaxTotalTime {
			return "", errors.New("timed out waiting for result")
		}

		result, err := c.GetStatus(ctx)
		if err != nil {
			log.Println("Error making request to GetStatus.")
			return "", err
		} else {
			if result == "completed" {
				return "completed", nil
			}
			if result == "error" {
				return "error", nil
			}
			if result == "pending" {
				log.Println("pending")
			}
		}

		// Wait for 'interval' before retrying or handle context cancellation
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		time.Sleep(interval)

		// Increase interval using exponential backoff until it hits the cap MaxWait
		interval = time.Duration(float64(interval) * c.BackoffFactor)
		if interval > c.MaxWait {
			interval = c.MaxWait
		}
	}
}