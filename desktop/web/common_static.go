package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:static/common
var staticCommonFS embed.FS

func handleCommonFileServer() http.Handler {
	sub, err := fs.Sub(staticCommonFS, "static/common")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}
