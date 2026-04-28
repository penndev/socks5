package web

import "net/http"

func Route(w http.ResponseWriter, r *http.Request) {
	router := newRouter()
	router.ServeHTTP(w, r)
}

func newRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handleRoot)
	router.HandleFunc("/subscribe", handleSubscribeRedirect)
	router.HandleFunc("/subscribe/logo.png", handleSubscribeLogo)
	router.HandleFunc("/subscribe/api/servers", handleSubscribeServers)
	router.HandleFunc("/subscribe/api/servers/import", handleSubscribeImportServers)
	router.HandleFunc("/subscribe/api/servers/export", handleSubscribeExportServers)
	router.HandleFunc("/subscribe/api/subscription/convert", handleSubscribeSubscriptionConvert)
	router.Handle("/subscribe/", http.StripPrefix("/subscribe/", handleSubscribeFileServer()))
	return router
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("/"))
}

func handleSubscribeRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/subscribe/", http.StatusFound)
}
