package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/penndev/gopkg/util"
)

// handleHTTPConnect 处理 CONNECT（常见于 HTTPS），建立双向隧道。
func (s *Server) handleHTTPConnect(client net.Conn, req *http.Request) error {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	if host == "" {
		fmt.Fprintf(client, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return fmt.Errorf("missing host")
	}
	// 回复 200 Connection Established
	fmt.Fprintf(client, "HTTP/1.1 200 Connection Established\r\n\r\n")
	// 建立隧道
	s.HandleConnect(client, "tcp", host)
	return nil
}

var hideHeader = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
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
	for _, key := range hideHeader {
		cc.Header.Del(key)
	}
	if err := cc.Write(local); err != nil {
		return err
	}
	util.Pipe(client, local)
	return nil
}

func (s *Server) ProxyHTTP(conn net.Conn) {
	br := bufio.NewReader(conn)
	req, err := http.ReadRequest(br)
	if err != nil {
		log.Println("read request failed: ", err)
		return
	}
	defer req.Body.Close()

	if req.Method == http.MethodConnect {
		if err := s.handleHTTPConnect(conn, req); err != nil {
			log.Println("connect failed: ", err)
		}
		return
	}

	isHttpProxy := req.URL.IsAbs() && strings.HasPrefix(req.URL.Scheme, "http")
	if isHttpProxy {
		if err := s.handleHTTPProxyForward(conn, req); err != nil {
			log.Println("http proxy forward failed: ", err)
		}
		return
	}
}
