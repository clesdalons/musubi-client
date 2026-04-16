package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/musubi-client/musubi-client/internal/cloud"
	"github.com/musubi-client/musubi-client/internal/config"
	"github.com/musubi-client/musubi-client/internal/storage"
	"github.com/musubi-client/musubi-client/internal/watcher"
)

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type App struct {
	ctx           context.Context
	watcher       *watcher.Watcher
	configManager *config.Manager
	azureClient   *cloud.AzureClient
}

func New() *App {
	return &App{
		configManager: config.NewManager(),
		azureClient:   cloud.NewClient(),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{Name: "Musubi", Version: "0.1.0"}
}

func (a *App) GetSettings() config.Config {
	cfg, _ := a.configManager.Load()
	return cfg
}

func (a *App) SaveSettings(path, uploader, campaign string) error {
	oldPath := a.GetSettings().SavePath
	cfg := config.Config{
		SavePath: path,
		Uploader: uploader,
		Campaign: campaign,
	}

	if err := a.configManager.Write(cfg); err != nil {
		return err
	}

	if oldPath != path && path != "" {
		return a.StartWatcher()
	}

	return nil
}

func (a *App) GetInitialPath() string {
	cfg, _ := a.configManager.Load()
	if cfg.SavePath != "" {
		return cfg.SavePath
	}

	defaultPath := a.configManager.DetectDefaultSavePath()
	if defaultPath == "" {
		return ""
	}

	cfg.SavePath = defaultPath
	_ = a.configManager.Write(cfg)
	return defaultPath
}

func (a *App) SelectFolder() string {
	folder, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select DOS2 Story Save Directory",
	})
	if err != nil || folder == "" {
		return ""
	}

	cfg := a.GetSettings()
	_ = a.SaveSettings(folder, cfg.Uploader, cfg.Campaign)
	return folder
}

func (a *App) DownloadLatestSave() string {
	if a.watcher != nil {
		a.watcher.Lock()
		log.Println("[Manager] Watcher locked for download")
	}

	defer func() {
		time.Sleep(3 * time.Second)
		if a.watcher != nil {
			a.watcher.Unlock()
			log.Println("[Manager] Watcher unlocked")
		}
	}()

	if err := a.downloadAndExtract(); err != nil {
		log.Printf("[Manager] Download failed: %v", err)
		return "Error: " + err.Error()
	}

	return "Success"
}

func (a *App) downloadAndExtract() error {
	cfg := a.GetSettings()
	if cfg.SavePath == "" {
		return fmt.Errorf("save path is not configured")
	}

	meta, tempPath, err := a.azureClient.DownloadLatestSave(cfg.Campaign)
	if err != nil {
		return err
	}
	defer os.Remove(tempPath)

	if err := storage.ExtractZip(tempPath, cfg.SavePath, meta.FileName); err != nil {
		return err
	}

	return nil
}

func (a *App) GetLocalSaveStatus() string {
	cfg := a.GetSettings()
	files, err := os.ReadDir(cfg.SavePath)
	if err != nil {
		return ""
	}

	newest := time.Time{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}

	if newest.IsZero() {
		return "Never"
	}

	return newest.Format(time.RFC3339)
}

func (a *App) GetCloudSaveStatus() (map[string]interface{}, error) {
	cfg := a.GetSettings()
	return a.azureClient.GetLatestSaveInfo(cfg.Campaign)
}

func (a *App) OpenFolder(path string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	_ = cmd.Run()
}

func (a *App) StartWatcher() error {
	cfg := a.GetSettings()
	if cfg.SavePath == "" {
		return nil
	}

	if a.watcher != nil {
		_ = a.watcher.Close()
	}

	w, err := watcher.NewWatcher(a.ctx, cfg.SavePath, a.handleSaveCreated)
	if err != nil {
		return err
	}

	a.watcher = w
	return nil
}

func (a *App) handleSaveCreated(savePath, saveName string) {
	if a.watcher != nil && a.watcher.IsLocked() {
		return
	}

	wailsRuntime.EventsEmit(a.ctx, "watcher:detected", saveName)
	go a.uploadSave(savePath)
}

func (a *App) uploadSave(savePath string) {
	zipPath, err := storage.ZipFolder(savePath)
	if err != nil {
		log.Printf("[Uploader] Zip error: %v", err)
		wailsRuntime.EventsEmit(a.ctx, "upload:error")
		return
	}
	defer os.Remove(zipPath)

	if err := a.azureClient.UploadSave(zipPath, a.GetSettings()); err != nil {
		log.Printf("[Uploader] Sync failed: %v", err)
		wailsRuntime.EventsEmit(a.ctx, "upload:error")
		return
	}

	log.Printf("[Uploader] Sync success: %s", savePath)
	wailsRuntime.EventsEmit(a.ctx, "upload:success")
}
