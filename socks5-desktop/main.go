package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

var app *application.App

func main() {
	storage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	app = application.New(application.Options{
		Name:        "socks5-desktop",
		Description: "Socks5 代理桌面应用",
		Icon:        appIcon,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(storage),
			application.NewService(&App{}),
			application.NewService(&Proxy{}),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "socks5-desktop",
		Width:            1000,
		Height:           800,
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: false,
		},
	})

	window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))

	var lastX, lastY int
	showAtLastPosition := func() {
		if lastX != 0 && lastY != 0 {
			window.SetPosition(lastX, lastY)
		}
		window.Show()
	}

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		lastX, lastY = window.Position()
		window.Hide()
		e.Cancel()
	})

	// 设置任务托管栏功能
	systray := app.SystemTray.New()
	systray.SetIcon(appIcon)
	systrayMenu := app.NewMenu()
	systrayMenu.Add("显示主窗口").OnClick(func(ctx *application.Context) {
		showAtLastPosition()
	})
	systrayMenu.Add("退出").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(systrayMenu)
	systray.OnClick(func() {
		if window.IsVisible() {
			lastX, lastY = window.Position()
			window.Hide()
		} else {
			showAtLastPosition()
		}
	})
	if err := app.Run(); err != nil {
		panic(err)
	}
}
