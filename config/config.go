package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	AuthToken string `json:"auth_token"`
	DBPath    string `json:"db_path,omitempty"`
}

// LoadAuthToken returns the auth token from (in order):
// 1. The provided flag value
// 2. THG_AUTH_TOKEN env var
// 3. ~/.config/thg/config.json
func LoadAuthToken(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if v := os.Getenv("THG_AUTH_TOKEN"); v != "" {
		return v, nil
	}

	cfg, err := loadConfigFile()
	if err == nil && cfg.AuthToken != "" {
		return cfg.AuthToken, nil
	}

	return "", fmt.Errorf(`auth token required for this command

Set it one of these ways:
  1. Flag:    --auth-token <token>
  2. Env var: export THG_AUTH_TOKEN=<token>
  3. Config:  echo '{"auth_token":"<token>"}' > ~/.config/thg/config.json

Find your token in Things → Settings → General → Enable Things URLs → Manage`)
}

func loadConfigFile() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	path := filepath.Join(home, ".config", "thg", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
