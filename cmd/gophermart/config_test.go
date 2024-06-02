package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoffIntervalConfig(t *testing.T) {
	testTable := []struct {
		name                     string
		envVars                  map[string]string
		expectedBackoffIntervals []time.Duration
		expectedBackoffRetries   bool
	}{
		{
			name:                     "Default configuration",
			envVars:                  map[string]string{},
			expectedBackoffIntervals: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
			expectedBackoffRetries:   true,
		},
		{
			name: "Custom backoff intervals",
			envVars: map[string]string{
				"BACKOFF_INTERVALS": "1s,2s,3s",
			},
			expectedBackoffIntervals: []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
			expectedBackoffRetries:   true,
		},
		{
			name: "Disable backoff retries",
			envVars: map[string]string{
				"BACKOFF_RETRIES": "false",
			},
			expectedBackoffIntervals: nil,
			expectedBackoffRetries:   false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = nil
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := NewConfig()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBackoffIntervals, cfg.Postgres.BackoffIntervals)
			assert.Equal(t, tt.expectedBackoffRetries, cfg.Postgres.BackoffRetries)

			for k := range tt.envVars {
				t.Setenv(k, "")
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	cfg, err := NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", cfg.ListenAddress)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, "http://localhost:8001", cfg.ActualSystemAddress)
	assert.Equal(t, 5, cfg.AccuralCheckerLimit)
	assert.Equal(t, "qwerty123", cfg.JWT.Secret)
	assert.Equal(t, 1*time.Hour, cfg.JWT.ExpiredDuration)
}

func TestEnvOverrides(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:9090")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://example.com")
	os.Setenv("ACCRUAL_CHECKER_LIMIT", "10")
	os.Setenv("JWT_SECRET", "secretKey")
	os.Setenv("JWT_EXPIRED", "2h")

	cfg, err := NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:9090", cfg.ListenAddress)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "http://example.com", cfg.ActualSystemAddress)
	assert.Equal(t, 10, cfg.AccuralCheckerLimit)
	assert.Equal(t, "secretKey", cfg.JWT.Secret)
	assert.Equal(t, 2*time.Hour, cfg.JWT.ExpiredDuration)

	// Cleanup
	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
	os.Unsetenv("ACCRUAL_CHECKER_LIMIT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRED")
}
