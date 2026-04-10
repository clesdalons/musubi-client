package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// AppInfo holds application metadata for the frontend
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

// GetAppInfo returns centralized metadata to maintain DRY principle
func (a *App) GetAppInfo() AppInfo {
	return AppInfo{
		Name:    "Musubi",
		Version: "0.1.0",
	}
}

// GetInitialPath returns the prioritized save directory (Saved Config > Auto-detect)
func (a *App) GetInitialPath() string {
	if saved := a.LoadConfig(); saved != "" {
		return saved
	}

	if auto := a.detectDefaultPath(); auto != "" {
		_ = a.SaveConfig(auto)
		return auto
	}

	return ""
}

// SelectFolder opens a native OS dialog and restarts the monitoring service if successful
func (a *App) SelectFolder() string {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select DOS2 Story Save Directory",
	})

	if err != nil || folder == "" {
		return ""
	}

	_ = a.SaveConfig(folder)
	a.StartWatcher()

	return folder
}
