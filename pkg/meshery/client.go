package meshery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client represents a client for Meshery API
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// SnapshotImportResponse represents the response from Meshery after importing a snapshot
type SnapshotImportResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewClient creates a new Meshery client
func NewClient(baseURL, token string) (*Client, error) {
	// Validate URL
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Meshery URL: %w", err)
	}

	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// ImportSnapshot imports a snapshot file to Meshery
func (c *Client) ImportSnapshot(ctx context.Context, filePath string) error {
	// Read the snapshot file
	snapshotData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %w", err)
	}

	// Create the request
	endpoint := fmt.Sprintf("%s/api/meshsync/snapshot/import", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(snapshotData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Meshery: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-success responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Meshery API returned non-success status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response SnapshotImportResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check response status
	if response.Status != "success" {
		return fmt.Errorf("import failed: %s", response.Message)
	}

	return nil
} 