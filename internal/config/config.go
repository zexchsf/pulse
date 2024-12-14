package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	Expiration    time.Duration
}

// loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "pulse_dev"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "pulse"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_SECRET", ""),
			RefreshSecret: getEnv("REFRESH_SECRET", ""),
			Expiration:    time.Duration(getEnvAsInt("JWT_EXPIRATION", 86400)) * time.Second,
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := strings.ToLower(getEnv(key, strconv.FormatBool(defaultValue)))

	switch valueStr {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// checks for required configuration values
func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if config.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if config.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if config.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if config.JWT.RefreshSecret == "" {
		return fmt.Errorf("REFRESH_SECRET is required")
	}

	return nil
}

// check if env is production
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Env) == "production"
}

// check if env is development
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Env) == "development"
}

// validate configs
func (c *Config) Validate() error {
	return validateConfig(c)
}
