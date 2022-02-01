package types

import (
	"fmt"
	"time"
)

// Config is used to hold application configuration values.
type Config struct {
	Port                int `mapstructure:"port"`
	KeyLength           int `mapstructure:"key_length"`
	CertificateLifeDays int `mapstructure:"cert_life_days"`
}

// GetCertificateDuration returns the cert lifetime duration.
func (c *Config) GetCertificateDuration() time.Duration {
	return time.Hour * 24 * time.Duration(c.CertificateLifeDays)
}

// String returns a string representation of config.
func (c Config) String() string {
	return fmt.Sprintf("port=%d key-length=%d cert-life=%v", c.Port, c.KeyLength, c.GetCertificateDuration())
}
