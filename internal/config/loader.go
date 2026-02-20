package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		path = "helixdb.config.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
