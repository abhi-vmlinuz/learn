package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	configDir  = ".config/learn"
	configFile = "config.toml"
)

type Config struct {
	Repo     RepoConfig     `toml:"repo"`
	Defaults DefaultsConfig `toml:"defaults"`
}

type RepoConfig struct {
	Root string `toml:"root"`
}

type DefaultsConfig struct {
	Viewer   string `toml:"viewer"`
	Category string `toml:"category"`
}

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, configDir, configFile)
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir)
}

func TemplatesDir() string {
	return filepath.Join(ConfigDir(), "templates")
}

func Load() (*Config, error) {
	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("not initialized: config not found at %s\nRun 'learn init' first", path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Repo.Root == "" {
		return nil, fmt.Errorf("repo root not set in config")
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	path := ConfigPath()
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(cfg)
}
