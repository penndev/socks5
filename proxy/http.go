package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/penndev/gopkg/util"
)

// handleHTTPConnect 处理 CONNECT（常见于 HTTPS），建立双向隧道。
// client 必须为已 Hijack 的连接，否则 net/http 在 handler 返回后可能继续写回包，破坏 TLS。
func (s *Server) handleHTTPConnect(client net.Conn, req *http.Request) error {
	host := req.Host
	if host == "" {
		fmt.Fprintf(client, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return fmt.Errorf("missing host")
	}
	fmt.Fprintf(client, "HTTP/1.1 200 Connection Established\r\n\r\n")
	s.HandleConnect(client, "tcp", host)
	return nil
}

// handleHTTPProxyForward 处理带绝对 URL 的 HTTP 代理请求（如 GET http://host/path）。
func (s *Server) handleHTTPProxyForward(client net.Conn, req *http.Request) error {
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}
	addr := net.JoinHostPort(req.URL.Hostname(), port)

	remote, local := net.Pipe()
	go func() {
		s.HandleConnect(remote, "tcp", addr)
	}()

	// 复写请求
	cc := req.Clone(req.Context())
	cc.RequestURI = ""
	cc.URL = &url.URL{
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}
	if cc.URL.Path == "" {
		cc.URL.Path = "/"
	}
	cc.Host = req.URL.Host
	if cc.URL.Path == "" {
		cc.URL.Path = "/"
	}
	for _, key := range HttpProxyHeaders {
		cc.Header.Del(key)
	}
	if err := cc.Write(local); err != nil {
		return err
	}
	util.Pipe(client, local)
	return nil
}

func (s *Server) ProxyHTTP(conn net.Conn) {
	listener := &HttpSingleConnListener{conn: conn}
	http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "hijacking not supported", http.StatusInternalServerError)
				return
			}
			client, _, err := hijacker.Hijack()
			if err != nil {
				log.Println("hijack failed: ", err)
				return
			}
			if err := s.handleHTTPConnect(client, r); err != nil {
				log.Println("connect failed: ", err)
			}
			return
		}
		// 传统http代理
		if r.URL.IsAbs() && strings.HasPrefix(r.URL.Scheme, "http") {
			if err := s.handleHTTPProxyForward(conn, r); err != nil {
				log.Println("http proxy forward failed: ", err)
			}
			return
		}
		// httpweb请求
		http.NotFound(w, r)
	}))
}
