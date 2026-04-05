package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// 转发普通 HTTP(S) 代理请求时直连上游，不走环境变量里的 HTTP_PROXY，避免回环。
var upstreamHTTP = &http.Client{
	Transport: &http.Transport{
		Proxy: nil,
	},
}

func main1() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	// CONNECT 的 URL.Path 常为空；走 DefaultServeMux 会触发「补 /」的 301，进不到业务 handler。
	// 必须在 ServeMux 之前处理 CONNECT。
	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleConnect(w, r)
				return
			}
			mux.ServeHTTP(w, r)
		}),
	}
	log.Println("HTTP proxy on http://127.0.0.1:8080 (CONNECT + 绝对 URL 转发)")
	log.Fatal(srv.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// 客户端走代理访问 http(s):// 时，请求行是完整绝对 URL，例如 GET http://example.com/a HTTP/1.1
	if r.URL.IsAbs() && (r.URL.Scheme == "http" || r.URL.Scheme == "https") {
		proxyHTTP(w, r)
		return
	}
	log.Println(r.URL.String(), r.Method)
	w.Write([]byte("ok\n"))
}

func proxyHTTP(w http.ResponseWriter, r *http.Request) {
	out := r.Clone(r.Context())
	out.RequestURI = ""
	stripHopByHop(out.Header)

	resp, err := upstreamHTTP.Do(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

var hopByHopHeaders = []string{
	"Connection", "Proxy-Connection", "Keep-Alive", "Proxy-Authenticate",
	"Proxy-Authorization", "Te", "Trailer", "Transfer-Encoding", "Upgrade",
}

func stripHopByHop(h http.Header) {
	for _, k := range hopByHopHeaders {
		h.Del(k)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if host == "" {
		host = r.URL.Host
	}
	if host == "" {
		http.Error(w, "missing host", http.StatusBadRequest)
		return
	}

	remote, err := net.DialTimeout("tcp", host, 15*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		remote.Close()
		http.Error(w, "hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, bufrw, err := hijacker.Hijack()
	if err != nil {
		remote.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.WriteString(clientConn, "HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
		clientConn.Close()
		remote.Close()
		return
	}

	go func() {
		io.Copy(remote, bufrw)
		remote.Close()
	}()
	io.Copy(clientConn, remote)
	remote.Close()
	clientConn.Close()
}
