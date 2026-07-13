package config

import "log/slog"

type Config struct {
	Dir    string
	DBPath string
	CACert string
	CAKey  string

	Proxy   ProxyConfig   `toml:"proxy"`
	Logging LoggingConfig `toml:"logging"`
	Watch   WatchConfig   `toml:"watch"`
}

type ProxyConfig struct {
	Port int `toml:"port"`
}

type LoggingConfig struct {
	Level slog.Level `toml:"level"`
}

type WatchConfig struct {
	Limit int `toml:"limit"`
}
