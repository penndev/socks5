package main

import (
	"desktop/build"
	"desktop/internal"
	"desktop/lang"
	"desktop/proxy"
	"desktop/storage"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	store, err := storage.New()
	if err != nil {
		panic(err)
	}
	language, err := lang.New()
	if err != nil {
		panic(err)
	}
	app, err := internal.SetUPApp([]application.Service{
		application.NewService(store),
		application.NewService(language),
		application.NewService(&internal.AppConst{}),
		application.NewService(&proxy.Proxy{}),
		application.NewService(&proxy.ProxyPing{}),
	})
	if err != nil {
		panic(err)
	}

	// 设置app托管栏
	systrayMenu := app.NewMenu()
	systrayMenuShowMain := systrayMenu.Add(language.T("systray.showMainWindow"))
	systrayMenuQuit := systrayMenu.Add(language.T("systray.quit"))

	// 设置app托管栏功能
	systrayMenuShowMain.OnClick(func(ctx *application.Context) {
		internal.MainWindowSetPosition()
	})
	systrayMenuQuit.OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// 设置任务托管栏功能
	systray := app.SystemTray.New()
	systray.SetIcon(build.Icon)
	systray.SetMenu(systrayMenu)
	systray.OnClick(func() {
		internal.MainWindowSetPosition()
	})

	// 设置任务托管栏文案随 Lang 切换更新
	app.Event.On(internal.AppConfig.EventNameLocaleChanged, func(_ *application.CustomEvent) {
		systrayMenuShowMain.SetLabel(language.T("systray.showMainWindow"))
		systrayMenuQuit.SetLabel(language.T("systray.quit"))
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
