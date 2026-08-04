package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".promtrace")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	cfg := &Config{
		Dir:    dir,
		DBPath: filepath.Join(dir, "traces.db"),
		CACert: filepath.Join(dir, "ca.crt"),
		CAKey:  filepath.Join(dir, "ca.key"),
		Proxy:  ProxyConfig{Port: 9117},
		Watch:  WatchConfig{Limit: 20},
	}

	tomlPath := filepath.Join(dir, "config.toml")
	if _, err := os.Stat(tomlPath); err == nil {
		_, err = toml.DecodeFile(tomlPath, cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
