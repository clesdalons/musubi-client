package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all persistent user settings
type Config struct {
	SavePath string `json:"save_path"`
	Uploader string `json:"uploader"`
	Campaign string `json:"campaign"`
}

// getConfigPath returns the path to the local config file
func (a *App) getConfigPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".musubi")
	_ = os.MkdirAll(path, 0755)
	return filepath.Join(path, "config.json")
}

// LoadConfig reads the config file from disk
func (a *App) LoadConfig() Config {
	var cfg Config
	data, err := os.ReadFile(a.getConfigPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

// WriteConfig saves all settings to the config file
func (a *App) WriteConfig(path, uploader, campaign string) error {
	cfg := Config{
		SavePath: path,
		Uploader: uploader,
		Campaign: campaign,
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.getConfigPath(), data, 0644)
}

// detectDefaultPath attempts to locate the DOS2 Savegames folder
func (a *App) detectDefaultPath() string {
	home, _ := os.UserHomeDir()

	// Support for standard and OneDrive paths
	possibleBases := []string{
		filepath.Join(home, "Documents"),
		filepath.Join(home, "OneDrive", "Documents"),
	}

	for _, base := range possibleBases {
		storyPath := filepath.Join(base, "Larian Studios", "Divinity Original Sin 2 Definitive Edition", "Player Profiles")
		if _, err := os.Stat(storyPath); err == nil {
			profiles, _ := os.ReadDir(storyPath)
			for _, profile := range profiles {
				if profile.IsDir() {
					fullPath := filepath.Join(storyPath, profile.Name(), "Savegames", "Story")
					if _, err := os.Stat(fullPath); err == nil {
						return fullPath
					}
				}
			}
		}
	}
	return ""
}
