package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// remote := r.URL.Hostname() + ":" + r.URL.Port()
			// log.Println(remote)

			if r.Method == http.MethodConnect {
				log.Println("proxy HTTPS")
				return
			}

			isProxyHTTP := r.URL.IsAbs() && strings.HasPrefix(r.URL.Scheme, "http")
			if isProxyHTTP {
				log.Println("proxy HTTP")
				return
			}

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("forbidden\n"))
		}),
	}
	log.Fatal(srv.ListenAndServe())
}
