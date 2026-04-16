package watcher

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher     *fsnotify.Watcher
	ctx           context.Context
	onSaveCreated func(savePath, saveName string)
	locked        bool
}

func NewWatcher(ctx context.Context, path string, onSaveCreated func(savePath, saveName string)) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fsWatcher.Add(path); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	w := &Watcher{
		fsWatcher:     fsWatcher,
		ctx:           ctx,
		onSaveCreated: onSaveCreated,
	}
	go w.listen()
	return w, nil
}

func (w *Watcher) listen() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			if w.locked {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err != nil || !info.IsDir() {
					continue
				}

				saveName := filepath.Base(event.Name)
				log.Printf("[Watcher] New save detected: %s", saveName)
				go func(path, name string) {
					time.Sleep(3 * time.Second)
					w.onSaveCreated(path, name)
				}(event.Name, saveName)
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("[Watcher] error: %v", err)
			return
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *Watcher) Close() error {
	if w.fsWatcher == nil {
		return nil
	}
	return w.fsWatcher.Close()
}

func (w *Watcher) Lock() {
	w.locked = true
}

func (w *Watcher) Unlock() {
	w.locked = false
}

func (w *Watcher) IsLocked() bool {
	return w.locked
}
