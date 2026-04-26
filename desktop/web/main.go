package web

import "net/http"

func Route(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	})
	router.ServeHTTP(w, r)
}
