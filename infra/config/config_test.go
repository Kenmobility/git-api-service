package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/stretchr/testify/assert"
)

func setupEnv(envs map[string]string) {
	for key, value := range envs {
		os.Setenv(key, value)
	}
}

func clearEnv(keys []string) {
	for _, key := range keys {
		os.Unsetenv(key)
	}
}

func TestLoadConfigValid(t *testing.T) {
	// Set up environment variables for a valid configuration
	envs := map[string]string{
		"APP_ENV":                   "test",
		"GIT_HUB_TOKEN":             "test_token",
		"DATABASE_HOST":             "localhost",
		"DATABASE_PORT":             "5432",
		"DATABASE_USER":             "test_user",
		"DATABASE_PASSWORD":         "test_password",
		"DATABASE_NAME":             "test_db",
		"FETCH_INTERVAL":            "1h",
		"DEFAULT_START_DATE":        "2023-01-01T00:00:00Z",
		"DEFAULT_END_DATE":          "2024-01-01T00:00:00Z",
		"GIT_COMMIT_FETCH_PER_PAGE": "100",
		"GITHUB_API_BASE_URL":       "https://api.github.com",
		"DEFAULT_REPOSITORY":        "example/repo",
		"ADDRESS":                   "127.0.0.1",
		"PORT":                      "8080",
	}

	setupEnv(envs)
	defer clearEnv([]string{
		"APP_ENV", "GIT_HUB_TOKEN", "DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER",
		"DATABASE_PASSWORD", "DATABASE_NAME", "FETCH_INTERVAL", "DEFAULT_START_DATE",
		"DEFAULT_END_DATE", "GIT_COMMIT_FETCH_PER_PAGE", "GITHUB_API_BASE_URL",
		"DEFAULT_REPOSITORY", "ADDRESS", "PORT",
	})

	// Load the config
	cfg, err := config.LoadConfig("")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check if all values are correctly loaded
	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "test_token", cfg.GitHubToken)
	assert.Equal(t, "localhost", cfg.DatabaseHost)
	assert.Equal(t, "5432", cfg.DatabasePort)
	assert.Equal(t, "test_user", cfg.DatabaseUser)
	assert.Equal(t, "test_password", cfg.DatabasePassword)
	assert.Equal(t, "test_db", cfg.DatabaseName)
	assert.Equal(t, time.Hour, cfg.FetchInterval)
	assert.Equal(t, 100, cfg.GitCommitFetchPerPage)
	assert.Equal(t, "https://api.github.com", cfg.GitHubApiBaseURL)
	assert.Equal(t, "127.0.0.1", cfg.Address)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "example/repo", cfg.DefaultRepository)

	// Check start and end dates
	expectedStartDate, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	expectedEndDate, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	assert.Equal(t, expectedStartDate, cfg.DefaultStartDate)
	assert.Equal(t, expectedEndDate, cfg.DefaultEndDate)
}

func TestLoadConfigInvalid(t *testing.T) {
	// Set up environment variables with invalid values
	envs := map[string]string{
		"APP_ENV":                   "test",
		"FETCH_INTERVAL":            "invalid_duration",
		"DEFAULT_START_DATE":        "invalid_date",
		"GIT_COMMIT_FETCH_PER_PAGE": "invalid_int",
	}
	setupEnv(envs)
	defer clearEnv([]string{"APP_ENV", "FETCH_INTERVAL", "DEFAULT_START_DATE", "GIT_COMMIT_FETCH_PER_PAGE"})

	// Attempt to load the config, expecting errors
	cfg, err := config.LoadConfig("")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfigMissingEnvVars(t *testing.T) {
	// Set up environment variables with some required ones missing
	envs := map[string]string{
		"APP_ENV":        "test",
		"DATABASE_PORT":  "5432",
		"FETCH_INTERVAL": "2h",
	}
	setupEnv(envs)
	defer clearEnv([]string{"APP_ENV", "DATABASE_PORT", "FETCH_INTERVAL"})

	// Attempt to load the config, expecting validation error
	cfg, err := config.LoadConfig("")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfigDefaultValues(t *testing.T) {
	// Test for default values when some environment variables are not set
	envs := map[string]string{
		"APP_ENV":           "test",
		"DATABASE_HOST":     "localhost",
		"DATABASE_PORT":     "5432",
		"DATABASE_USER":     "test_user",
		"DATABASE_PASSWORD": "test_password",
		"DATABASE_NAME":     "test_db",
	}
	setupEnv(envs)
	defer clearEnv([]string{"APP_ENV", "DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_NAME"})

	// Load config with missing optional env vars
	cfg, err := config.LoadConfig("")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check if default values are applied
	assert.Equal(t, time.Hour, cfg.FetchInterval)
	assert.Equal(t, "chromium/chromium", cfg.DefaultRepository)
}
