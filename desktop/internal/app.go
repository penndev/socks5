package internal

import (
	"desktop/build"
	"desktop/frontend"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var App *application.App

func SetUPApp(services []application.Service) (*application.App, error) {

	App = application.New(application.Options{
		Name:        "Prism",
		Description: "Prism APP",
		Icon:        build.Icon,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(frontend.Assets),
		},
		Services: services,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	SetUPMainWindow()
	return App, nil
}
