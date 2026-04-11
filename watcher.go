package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	IsLocked  bool
}

func (a *App) StartWatcher() {
	if a.watcher != nil {
		a.watcher.fsWatcher.Close()
	}

	cfg := a.LoadConfig()
	if cfg.SavePath == "" {
		return
	}

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("[Watcher] Error: %v", err)
		return
	}

	a.watcher = &Watcher{fsWatcher: fs}
	go a.listenToEvents()

	_ = a.watcher.fsWatcher.Add(cfg.SavePath)
	log.Printf("[Watcher] Monitoring: %s", cfg.SavePath)
}

func (a *App) listenToEvents() {
	for {
		select {
		case event, ok := <-a.watcher.fsWatcher.Events:
			if !ok {
				return
			}

			// Prevent upload while writing the downloaded file from musubi-api
			if a.watcher.IsLocked {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					saveName := filepath.Base(event.Name)
					savePath := event.Name

					runtime.EventsEmit(a.ctx, "watcher:detected", saveName)

					go func(p, n string) {
						time.Sleep(3 * time.Second) // Safety delay for game I/O

						zipPath, err := a.ZipFolder(p)
						if err != nil {
							log.Printf("[Uploader] Zip Error: %v", err)
							return
						}

						if err := a.UploadToAzure(zipPath); err != nil {
							log.Printf("[Uploader] Sync Failed: %v", err)
							runtime.EventsEmit(a.ctx, "upload:error")
						} else {
							log.Printf("[Uploader] Sync Success: %s", n)
							runtime.EventsEmit(a.ctx, "upload:success")
						}
					}(savePath, saveName)
				}
			}
		case <-a.watcher.fsWatcher.Errors:
			return
		}
	}
}
