package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

const (
	InitFlagFilename = "init"
)

type ProjectInitFlag struct {
	Initialized bool `json:"initialized"`
}

// TODO: we need to remove the global config instance keeping it now just until everything is migrated
var instance atomic.Pointer[Config]

func Init(workingDir string, debug bool) (*Config, error) {
	cfg, err := Load(workingDir, debug)
	if err != nil {
		return nil, err
	}
	instance.Store(cfg)
	return instance.Load(), nil
}

func Get() *Config {
	cfg := instance.Load()
	return cfg
}

func ProjectNeedsInitialization() (bool, error) {
	cfg := Get()
	if cfg == nil {
		return false, fmt.Errorf("config not loaded")
	}

	flagFilePath := filepath.Join(cfg.Options.DataDirectory, InitFlagFilename)

	_, err := os.Stat(flagFilePath)
	if err == nil {
		return false, nil
	}

	if !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to check init flag file: %w", err)
	}

	openpilotExists, err := openpilotMdExists(cfg.WorkingDir())
	if err != nil {
		return false, fmt.Errorf("failed to check for OPENPILOT.md files: %w", err)
	}
	if openpilotExists {
		return false, nil
	}

	return true, nil
}

func openpilotMdExists(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := strings.ToLower(entry.Name())
		if name == "openpilot.md" || name == "openpilot.local.md" {
			return true, nil
		}
	}

	return false, nil
}

func MarkProjectInitialized() error {
	cfg := Get()
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}
	flagFilePath := filepath.Join(cfg.Options.DataDirectory, InitFlagFilename)

	file, err := os.Create(flagFilePath)
	if err != nil {
		return fmt.Errorf("failed to create init flag file: %w", err)
	}
	defer file.Close()

	return nil
}

func HasInitialDataConfig() bool {
	cfgPath := GlobalConfigData()
	if _, err := os.Stat(cfgPath); err != nil {
		return false
	}
	return Get().IsConfigured()
}
