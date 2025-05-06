package settings_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/settings"
)

const testConfigFile = "config.json"

func cleanupConfigFile(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(testConfigFile); err == nil {
		if err := os.Remove(testConfigFile); err != nil {
			t.Fatalf("Failed to clean up config file: %v", err)
		}
	}
}

func TestLoadConfigCreatesDefaultIfMissing(t *testing.T) {
	cleanupConfigFile(t)
	defer cleanupConfigFile(t)

	cfg, err := settings.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error loading default config, got %v", err)
	}

	if len(cfg.CommitTypes) == 0 {
		t.Errorf("Expected default commit types to be present")
	}
	if cfg.DefaultCommitType == "" {
		t.Errorf("Expected default commit type to be set")
	}

	if _, err := os.Stat(testConfigFile); os.IsNotExist(err) {
		t.Errorf("Expected config.json to be created, but it doesn't exist")
	}
}

func TestLoadConfigReadsExistingFile(t *testing.T) {
	defer cleanupConfigFile(t)

	// Prepare a valid test config file
	testCfg := settings.Config{
		CommitTypes:          []string{"test"},
		CommitFormat:         "$type: $summary",
		DefaultVersion:       "1.0",
		DefaultCommitType:    "test",
		DefaultJiraReference: "TEST-123",
	}
	data, _ := json.MarshalIndent(testCfg, "", "  ")
	if err := os.WriteFile(testConfigFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := settings.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error loading config, got %v", err)
	}

	if cfg.DefaultCommitType != "test" {
		t.Errorf("Expected DefaultCommitType to be 'test', got '%s'", cfg.DefaultCommitType)
	}
}
