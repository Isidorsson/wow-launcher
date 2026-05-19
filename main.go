package main

import (
	"embed"

	"wow-launcher/internal/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed config.toml
var configToml []byte

func main() {
	config.SetEmbedded(configToml)

	app := NewApp()

	cfg, err := config.Load()
	title := "WoW Launcher"
	if err == nil {
		title = cfg.Branding.WindowTitle
	}

	err = wails.Run(&options.App{
		Title:  title,
		Width:  1100,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 22, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
