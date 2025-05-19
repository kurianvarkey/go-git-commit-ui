package settings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed default_config.json
var embeddedDefaultConfigData []byte // Embed the default config file content as a byte slice

// Config represents the structure of our configuration file.
type Config struct {
	CommitTypes          []string `json:"commit_types"` // Alias for git_commit_types
	CommitFormat         string   `json:"commit_format"`
	DefaultVersion       string   `json:"default_version"`
	DefaultCommitType    string   `json:"default_commit_type"`
	DefaultJiraReference string   `json:"default_jira_reference"`
}

const configFileName = "git-commit-ui-config.json"

// LoadConfig attempts to load a Config object from disk. If the file does not exist, it will be created
// with default values. If the file exists, it will be read from disk and deserialized into a Config
// object. If errors occur during the loading process, an error will be returned.
func LoadConfig() (*Config, error) {
	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		return createDefaultConfig()
	} else if err != nil {
		return nil, fmt.Errorf("error checking %s status: %w", configFileName, err)
	}

	configFileData, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s from disk: %w", configFileName, err)
	}

	var config Config
	err = json.Unmarshal(configFileData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s from disk: %w", configFileName, err)
	}

	return &config, nil
}

// createDefaultConfig creates a default configuration file from the embedded data,
// and returns the deserialized configuration object. If any errors occur while
// writing the file or unmarshaling the embedded data, an error is returned.
func createDefaultConfig() (*Config, error) {
	err := os.WriteFile(configFileName, embeddedDefaultConfigData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write default config file from embedded data: %w", err)
	}

	var config Config
	err = json.Unmarshal(embeddedDefaultConfigData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedded default config: %w", err)
	}

	return &config, nil
}
