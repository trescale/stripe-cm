package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	keychainService = "stripe-cm"
	keychainAccount = "stripe-cm"
)

type Config struct {
	StripeAPIKey string `yaml:"stripe_api_key"`
	Storage      string `yaml:"storage"` // "file" or "keychain"
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".stripe-cm"
	}
	return filepath.Join(home, ".stripe-cm")
}

func Load() (*Config, error) {
	cfgPath := filepath.Join(ConfigDir(), "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", cfgPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Storage == "keychain" {
		key, err := getKeychainKey()
		if err != nil {
			return nil, fmt.Errorf("read keychain: %w", err)
		}
		cfg.StripeAPIKey = key
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	toSave := *cfg
	if cfg.Storage == "keychain" {
		if err := setKeychainKey(cfg.StripeAPIKey); err != nil {
			return fmt.Errorf("store in keychain: %w", err)
		}
		toSave.StripeAPIKey = ""
	}

	data, err := yaml.Marshal(&toSave)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func getKeychainKey() (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", keychainService, "-w")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("keychain lookup failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func setKeychainKey(key string) error {
	cmd := exec.Command("security", "add-generic-password",
		"-s", keychainService,
		"-a", keychainAccount,
		"-w", key,
		"-U",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("keychain store failed: %s: %w", string(out), err)
	}
	return nil
}
