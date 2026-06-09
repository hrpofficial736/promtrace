package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Dir    string
	DBPath string
	CACert string
	CAKey  string
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".promtrace")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	return &Config{
		Dir:    dir,
		DBPath: filepath.Join(dir, "traces.db"),
		CACert: filepath.Join(dir, "ca.crt"),
		CAKey:  filepath.Join(dir, "ca.key"),
	}, nil
}
