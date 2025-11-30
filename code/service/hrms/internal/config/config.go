package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	IDGen      IDGenConfig
	Boundary   BoundaryConfig
	Individual IndividualConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port        string
	ContextPath string
}

// DatabaseConfig holds database related configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// IDGenConfig holds configuration for the ID generation service
type IDGenConfig struct {
	Host      string `mapstructure:"host"`
	Path      string `mapstructure:"path"`
	Enabled   bool   `mapstructure:"enabled"`
	IDGenName string `mapstructure:"idgen_name"`
}

type BoundaryConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type IndividualConfig struct {
	Host    string `mapstructure:"host"`
	Path    string `mapstructure:"path"`
	Timeout bool   `mapstructure:"timeout"`
}

// ServiceConfig holds configuration for an external service
type ServiceConfig struct {
	Host           string `mapstructure:"host"`
	CreateEndpoint string `mapstructure:"create_endpoint"`
	UpdateEndpoint string `mapstructure:"update_endpoint"`
	SearchEndpoint string `mapstructure:"search_endpoint"`
	DeleteEndpoint string `mapstructure:"delete_endpoint"`
	Timeout        int    `mapstructure:"timeout"`
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			ContextPath: getEnv("SERVER_CONTEXT_PATH", "/hrms"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "1234"),
			DBName:   getEnv("DB_NAME", "hrms_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		IDGen: IDGenConfig{
			Host:      getEnv("IDGEN_HOST", "http://localhost:8100"),
			Path:      getEnv("IDGEN_PATH", "/idgen/v1/generate"),
			Enabled:   getEnvAsBool("IDGEN_ENABLED", true),
			IDGenName: getEnv("IDGEN_NAME", "hrms.idgen"),
		},
		Boundary: BoundaryConfig{
			BaseURL: getEnv("BOUNDARY_HOST", "http://localhost:8095"),
		},
		Individual: IndividualConfig{
			Host:    getEnv("INDIVIDUAL_HOST", "http://localhost:8086"),
			Path:    getEnv("INDIVIDUAL_PATH", "/individual/v1"),
			Timeout: getEnvAsBool("INDIVIDUAL_TIMEOUT", true),
		},
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Helper functions matching style of other services
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		lowered := strings.ToLower(val)
		return lowered == "1" || lowered == "true" || lowered == "yes"
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	if val := os.Getenv(key); val != "" {
		parts := strings.Split(val, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			s := strings.TrimSpace(p)
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return defaultVal
}
