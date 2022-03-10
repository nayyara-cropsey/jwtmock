package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// Environment variable overrides
	portEnv     = "port"
	keyLenEnv   = "key_length"
	certLifeEnv = "cert_life_days"
	logLevelEnv = "log_level"

	envPrefix = "JWT_MOCK"
)

// Config is used to hold application configuration values.
type Config struct {
	Port                int    `yaml:"port"`
	KeyLength           int    `yaml:"key_length"`
	CertificateLifeDays int    `yaml:"cert_life_days"`
	LogLevel            string `yaml:"log_level"`
}

// GetCertificateDuration returns the cert lifetime duration.
func (c *Config) GetCertificateDuration() time.Duration {
	return time.Hour * 24 * time.Duration(c.CertificateLifeDays)
}

// String returns a string representation of config.
func (c Config) String() string {
	return fmt.Sprintf("port=%d key-length=%d cert-life=%v", c.Port, c.KeyLength, c.GetCertificateDuration())
}

// LoadConfig reads the given YAML file and loads config from it
func LoadConfig(yamlFileName string) (*Config, error) {
	contents, err := os.ReadFile(yamlFileName)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(contents, &cfg); err != nil {
		return nil, fmt.Errorf("yaml parse: %w", err)
	}

	if val, ok := getEnvVarInt(portEnv); ok {
		cfg.Port = val
	}

	if val, ok := getEnvVarInt(keyLenEnv); ok {
		cfg.KeyLength = val
	}

	if val, ok := getEnvVarInt(certLifeEnv); ok {
		cfg.CertificateLifeDays = val
	}

	if val, ok := getEnvVarStr(logLevelEnv); ok {
		cfg.LogLevel = val
	}

	return &cfg, nil
}

func getEnvVarInt(s string) (int, bool) {
	varName := fmt.Sprintf("%v_%v", envPrefix, strings.ToUpper(s))
	varStr := os.Getenv(varName)
	if varStr == "" {
		return 0, false
	}

	parseInt, err := strconv.Atoi(varStr)
	if err != nil {
		return 0, false
	}

	return parseInt, true
}
func getEnvVarStr(s string) (string, bool) {
	varName := fmt.Sprintf("%v_%v", envPrefix, strings.ToUpper(s))
	varStr := os.Getenv(varName)
	var present bool
	if varStr != "" {
		present = true
	}

	return varStr, present
}


