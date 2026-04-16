package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SavePath string `json:"save_path"`
	Uploader string `json:"uploader"`
	Campaign string `json:"campaign"`
}

func (c Config) IsComplete() bool {
	return c.SavePath != "" && c.Uploader != "" && c.Campaign != ""
}

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) configFilePath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".musubi")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "config.json")
}

func (m *Manager) Load() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(m.configFilePath())
	if err != nil {
		return cfg, nil
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (m *Manager) Write(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configFilePath(), data, 0o644)
}

func (m *Manager) DetectDefaultSavePath() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}

	possibleBaseDirs := []string{
		filepath.Join(home, "Documents"),
		filepath.Join(home, "OneDrive", "Documents"),
	}

	for _, base := range possibleBaseDirs {
		profilesBase := filepath.Join(base, "Larian Studios", "Divinity Original Sin 2 Definitive Edition", "Player Profiles")
		if _, err := os.Stat(profilesBase); err != nil {
			continue
		}

		profiles, err := os.ReadDir(profilesBase)
		if err != nil {
			continue
		}

		for _, profile := range profiles {
			if !profile.IsDir() {
				continue
			}

			candidate := filepath.Join(profilesBase, profile.Name(), "Savegames", "Story")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	return ""
}
