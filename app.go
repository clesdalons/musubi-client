package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type App struct {
	ctx     context.Context
	watcher *Watcher
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{Name: "Musubi", Version: "0.1.0"}
}

func (a *App) GetSettings() Config {
	return a.LoadConfig()
}

func (a *App) SaveSettings(path, uploader, campaign string) {
	oldPath := a.LoadConfig().SavePath
	_ = a.WriteConfig(path, uploader, campaign)

	// Restart watcher if the folder path has changed
	if oldPath != path && path != "" {
		a.StartWatcher()
	}
}

func (a *App) GetInitialPath() string {
	cfg := a.LoadConfig()
	if cfg.SavePath != "" {
		return cfg.SavePath
	}
	if auto := a.detectDefaultPath(); auto != "" {
		_ = a.WriteConfig(auto, cfg.Uploader, cfg.Campaign)
		return auto
	}
	return ""
}

func (a *App) SelectFolder() string {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select DOS2 Story Save Directory",
	})
	if err != nil || folder == "" {
		return ""
	}
	cfg := a.LoadConfig()
	a.SaveSettings(folder, cfg.Uploader, cfg.Campaign)
	return folder
}

func (a *App) DownloadLatestSave() string {
	// 1. Lock the watcher
	if a.watcher != nil {
		a.watcher.IsLocked = true
		log.Println("[Manager] Watcher locked for download")
	}

	// 2. Perform download
	err := a.DownloadAndExtract()

	// 3. Unlock with a small safety buffer
	go func() {
		time.Sleep(3 * time.Second)
		if a.watcher != nil {
			a.watcher.IsLocked = false
			log.Println("[Manager] Watcher unlocked")
		}
	}()

	if err != nil {
		log.Printf("[Manager] Download failed: %v", err)
		return "Error: " + err.Error()
	}

	return "Success"
}

// GetLocalSaveStatus returns the timestamp of the newest folder in the save path
func (a *App) GetLocalSaveStatus() string {
	cfg := a.LoadConfig()
	files, err := os.ReadDir(cfg.SavePath)
	if err != nil {
		return ""
	}

	var newestTime time.Time
	for _, f := range files {
		if f.IsDir() {
			info, _ := f.Info()
			if info.ModTime().After(newestTime) {
				newestTime = info.ModTime()
			}
		}
	}
	if newestTime.IsZero() {
		return "Never"
	}
	return newestTime.Format(time.RFC3339)
}

// GetCloudSaveStatus calls your new GetLatestSaveInfo endpoint
func (a *App) GetCloudSaveStatus() (map[string]interface{}, error) {
	cfg := a.LoadConfig()
	url := fmt.Sprintf("https://musubi.azurewebsites.net/api/GetLatestSaveInfo?campaignId=%s", cfg.Campaign)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
