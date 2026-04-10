package main

import (
	"context"
	"log"
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
