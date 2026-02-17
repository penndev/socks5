//go:build !windows

package main

type App struct{}

func (a *App) SetLaunchAtStartup(enabled bool) error {
	return nil
}

func (a *App) IsLaunchAtStartup() (bool, error) {
	return false, nil
}
