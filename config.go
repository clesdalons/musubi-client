package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SavePath string `json:"save_path"`
}

// getConfigPath returns the persistent storage location for Musubi settings
func (a *App) getConfigPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".musubi")
	_ = os.MkdirAll(path, 0755)
	return filepath.Join(path, "config.json")
}

// SaveConfig persists the chosen path to disk
func (a *App) SaveConfig(path string) error {
	data, err := json.MarshalIndent(Config{SavePath: path}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.getConfigPath(), data, 0644)
}

// LoadConfig retrieves the saved path from disk
func (a *App) LoadConfig() string {
	data, err := os.ReadFile(a.getConfigPath())
	if err != nil {
		return ""
	}
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	return cfg.SavePath
}

// detectDefaultPath attempts to locate the DOS2 Savegames folder automatically.
// It iterates through Player Profiles to find a valid 'Story' save directory.
func (a *App) detectDefaultPath() string {
	home, _ := os.UserHomeDir()

	// Base path for Definitive Edition
	basePath := filepath.Join(home, "Documents", "Larian Studios", "Divinity Original Sin 2 Definitive Edition", "PlayerProfiles")

	if _, err := os.Stat(basePath); err != nil {
		return ""
	}

	// List all profile folders (e.g., 'Public', 'MyPlayerName')
	profiles, err := os.ReadDir(basePath)
	if err != nil || len(profiles) == 0 {
		return ""
	}

	for _, profile := range profiles {
		if profile.IsDir() {
			// Check if this profile has a 'Story' save folder
			storyPath := filepath.Join(basePath, profile.Name(), "Savegames", "Story")
			if _, err := os.Stat(storyPath); err == nil {
				return storyPath
			}
		}
	}
	return ""
}

// PurgeConfig deletes the configuration file to reset application settings
func (a *App) PurgeConfig() error {
	configPath := a.getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		return os.Remove(configPath)
	}
	return nil
}
