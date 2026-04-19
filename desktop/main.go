package main

import (
	"embed"

	"desktop/internal"
	"desktop/lang"
	"desktop/proxy"
	"desktop/storage"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	langSvc, err := lang.New()
	if err != nil {
		panic(err)
	}

	internal.App = application.New(application.Options{
		Name:        "Prism",
		Description: "Prism代理桌面应用",
		Icon:        appIcon,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(store),
			application.NewService(&internal.AppConst{}),
			application.NewService(&proxy.Proxy{}),
			application.NewService(&proxy.ProxyPing{}),
			application.NewService(langSvc),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	app := internal.App
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Prism",
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
		window.Focus()
	}

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		lastX, lastY = window.Position()
		window.Hide()
		e.Cancel()
	})

	// 设置任务托管栏功能（文案随 Lang 切换更新）
	systray := app.SystemTray.New()
	systray.SetIcon(appIcon)
	systrayMenu := app.NewMenu()
	itemShowMain := systrayMenu.Add(langSvc.T("systray.showMainWindow")).OnClick(func(ctx *application.Context) {
		showAtLastPosition()
	})
	itemQuit := systrayMenu.Add(langSvc.T("systray.quit")).OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(systrayMenu)

	updateSystrayLabels := func() {
		itemShowMain.SetLabel(langSvc.T("systray.showMainWindow"))
		itemQuit.SetLabel(langSvc.T("systray.quit"))
	}
	app.Event.On(internal.AppConfig.EventNameLocaleChanged, func(_ *application.CustomEvent) {
		updateSystrayLabels()
	})
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
