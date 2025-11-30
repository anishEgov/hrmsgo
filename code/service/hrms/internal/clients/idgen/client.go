package idgen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Client handles ID generation service integration
type Client interface {
	// GenerateIDs generates formatted IDs with optional custom variables
	GenerateIDs(ctx context.Context, tenantID string, count int, customVars map[string]string) ([]string, error)
}

type client struct {
	host      string
	path      string
	enabled   bool
	idgenName string
}

// Config holds ID generation client configuration
type Config struct {
	Host      string
	Path      string
	Enabled   bool
	IDGenName string
}

// NewClient creates a new ID generation client
func NewClient(cfg Config) Client {
	return &client{
		host:      cfg.Host,
		path:      cfg.Path,
		enabled:   cfg.Enabled,
		idgenName: cfg.IDGenName,
	}
}

type idGenRequest struct {
	TemplateCode string            `json:"templateCode"`
	Variables    map[string]string `json:"variables"`
}

type idGenResponse struct {
	ID string `json:"id"`
}

// GenerateIDs generates formatted IDs from IDGen service
func (c *client) GenerateIDs(ctx context.Context, tenantID string, count int, customVars map[string]string) ([]string, error) {
	if !c.enabled {
		logrus.Warn("IDGen service is disabled, generating fallback IDs")
		return c.generateFallbackIDs(count), nil
	}

	url := c.host + c.path

	if customVars == nil {
		customVars = make(map[string]string)
	}

	// Ensure ORG is set from tenantID if not provided
	if _, exists := customVars["ORG"]; !exists && tenantID != "" {
		customVars["ORG"] = tenantID
	}

	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		payload := idGenRequest{
			TemplateCode: c.idgenName,
			Variables:    customVars,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal IDGen request: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		if tenantID != "" {
			req.Header.Set("X-Tenant-ID", tenantID)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Errorf("Failed to call IDGen service: %v. Using fallback ID.", err)
			ids = append(ids, c.generateFallbackIDs(1)[0])
			continue
		}

		func() {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				logrus.Errorf("IDGen service returned %d: %s. Using fallback ID.", resp.StatusCode, string(body))
				ids = append(ids, c.generateFallbackIDs(1)[0])
				return
			}

			var response idGenResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				logrus.Errorf("Failed to decode IDGen response: %v. Using fallback ID.", err)
				ids = append(ids, c.generateFallbackIDs(1)[0])
				return
			}

			if response.ID == "" {
				logrus.Error("IDGen response missing 'id'. Using fallback ID.")
				ids = append(ids, c.generateFallbackIDs(1)[0])
				return
			}

			ids = append(ids, response.ID)
		}()
	}

	return ids, nil
}

// generateFallbackIDs generates simple UUID-based IDs when IDGen service is unavailable
func (c *client) generateFallbackIDs(count int) []string {
	ids := make([]string, count)
	for i := 0; i < count; i++ {
		ids[i] = "EMP-" + uuid.New().String()[:8]
	}
	return ids
}
