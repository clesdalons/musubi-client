package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const (
	AppTitle  = "Musubi"
	MinWidth  = 450
	MinHeight = 740
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  AppTitle,
		Width:  MinWidth,
		Height: MinHeight,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 20, B: 25, A: 1},

		OnStartup: app.startup,

		// Wait for the UI to be fully mounted before starting background services
		OnDomReady: func(ctx context.Context) {
			app.StartWatcher()
		},

		// Expose the app instance to the frontend (Wails JS Bridge)
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
