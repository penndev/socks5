package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	storage, err := NewStorage()
	if err != nil {
		panic(err)
	}

	app := application.New(application.Options{
		Name:        "socks5-desktop",
		Description: "Socks5 代理桌面应用",
		Icon:        appIcon,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(storage),
			application.NewService(&App{}),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "socks5-desktop",
		Width:            400,
		Height:           800,
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: false,
		},
	})

	var lastX, lastY int
	var hasSavedPosition bool

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		lastX, lastY = window.Position()
		hasSavedPosition = true
		window.Hide()
		e.Cancel()
	})

	systray := app.SystemTray.New()
	systray.SetIcon(appIcon)
	systray.SetLabel("Socks5")
	systray.SetTooltip("Socks5 代理")

	showAtLastPosition := func() {
		if hasSavedPosition {
			window.SetPosition(lastX, lastY)
		}
		window.Show()
	}

	menu := app.NewMenu()
	menu.Add("显示主窗口").OnClick(func(ctx *application.Context) {
		showAtLastPosition()
	})
	menu.Add("退出").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)
	systray.OnClick(func() {
		if window.IsVisible() {
			lastX, lastY = window.Position()
			hasSavedPosition = true
			window.Hide()
		} else {
			showAtLastPosition()
		}
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
