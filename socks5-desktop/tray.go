package main

import (
	"context"
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

type Tray struct {
	appCtx context.Context
}

func NewTray() *Tray {
	return &Tray{}
}

func (t *Tray) Startup(ctx context.Context) {
	t.appCtx = ctx
	go systray.Run(t.onReady, t.onExit)
}

func (t *Tray) SetIconStatus() {
	systray.SetIcon(trayIcon)
}

func (t *Tray) onReady() {
	// systray.SetIcon(trayIcon)
	t.SetIconStatus()
	systray.SetTitle("Socks5")
	systray.SetTooltip("Socks5 代理")

	mShow := systray.AddMenuItem("显示主窗口", "显示主窗口")
	mQuit := systray.AddMenuItem("退出", "退出应用")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if t.appCtx != nil {
					runtime.WindowShow(t.appCtx)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				if t.appCtx != nil {
					runtime.Quit(t.appCtx)
				}
				return
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.appCtx = nil
}
