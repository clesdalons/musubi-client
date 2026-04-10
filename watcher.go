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

			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					saveName := filepath.Base(event.Name)
					savePath := event.Name

					runtime.EventsEmit(a.ctx, "new-save-event", saveName)

					go func(p, n string) {
						time.Sleep(3 * time.Second) // Safety delay for game I/O

						zipPath, err := a.ZipFolder(p)
						if err != nil {
							log.Printf("[Uploader] Zip Error: %v", err)
							return
						}

						if err := a.UploadToAzure(zipPath); err != nil {
							log.Printf("[Uploader] Sync Failed: %v", err)
						} else {
							log.Printf("[Uploader] Sync Success: %s", n)
						}
					}(savePath, saveName)
				}
			}
		case <-a.watcher.fsWatcher.Errors:
			return
		}
	}
}
