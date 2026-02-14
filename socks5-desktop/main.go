package main

import (
	"context"
	"embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 数据存储结构体实例
	storage, err := NewStorage()
	if err != nil {
		println("Error:", err.Error())
	}

	// 系统托盘结构体实例
	tray := NewTray()

	err = wails.Run(&options.App{
		Title:             "socks5-desktop",
		Width:             400,
		Height:            800,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			storage.Startup(ctx)
			tray.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			systray.Quit()
		},
		Bind: []any{
			storage,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
