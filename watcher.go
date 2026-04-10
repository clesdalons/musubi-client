package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
}

// StartWatcher initializes or restarts the background monitoring service
func (a *App) StartWatcher() {
	if a.watcher != nil {
		a.watcher.fsWatcher.Close()
	}

	path := a.LoadConfig()
	if path == "" {
		path = a.GetInitialPath()
	}

	if path == "" {
		log.Println("[Watcher] No valid path to monitor")
		return
	}

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("[Watcher] Creation error: %v", err)
		return
	}

	a.watcher = &Watcher{fsWatcher: fs}
	go a.listenToEvents()

	if err := a.watcher.fsWatcher.Add(path); err != nil {
		log.Printf("[Watcher] Error adding path %s: %v", path, err)
	} else {
		log.Printf("[Watcher] Started monitoring: %s", path)
	}
}

// listenToEvents handles incoming FS signals.
// In DOS2, a save is a folder created inside the 'Story' directory.
func (a *App) listenToEvents() {
	for {
		select {
		case event, ok := <-a.watcher.fsWatcher.Events:
			if !ok {
				return
			}

			// We monitor directory creation as DOS2 creates a new folder for each save
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)

				// Verify if the created item is indeed a directory (the save bundle)
				if err == nil && info.IsDir() {
					saveFolderName := filepath.Base(event.Name)
					log.Printf("[Watcher] New save detected: %s", saveFolderName)

					// Notify the React frontend via the Wails event bridge
					runtime.EventsEmit(a.ctx, "new-save-event", saveFolderName)
				}
			}

		case err, ok := <-a.watcher.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("[Watcher] Runtime error: %v", err)
		}
	}
}
