// internal/clients/boundary/client.go
package boundary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type Boundary struct {
	Code string `json:"code"`
	Name string `json:"name"`
	// Add other fields as needed
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    strings.TrimSuffix(baseURL, "/"),
	}
}

func (c *Client) SearchByCodes(ctx context.Context, tenantID string, codes []string) ([]Boundary, error) {
	if len(codes) == 0 {
		return []Boundary{}, nil
	}

	// Create request URL
	reqURL := fmt.Sprintf("%s/boundary/v1?codes=%s",
		c.baseURL,
		url.QueryEscape(strings.Join(codes, ",")))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("boundary service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("boundary service returned status: %d", resp.StatusCode)
	}

	var boundaries []Boundary
	if err := json.NewDecoder(resp.Body).Decode(&boundaries); err != nil {
		return nil, fmt.Errorf("failed to decode boundary service response: %w", err)
	}

	return boundaries, nil
}
