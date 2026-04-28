package web

import "net/http"

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("/"))
}
