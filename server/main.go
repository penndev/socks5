package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/penndev/gopkg/util"
	"github.com/penndev/socks5/core/transport"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println("listen failed: ", err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept failed: ", err)
			continue
		}
		go handleConn(NewConn(conn))
	}
}

func handleConn(conn *Conn) {
	// 判断协议类型如果是05则代理socks5，否则代理http
	buf, err := conn.Peek(1)
	if err != nil {
		log.Println("read failed: ", err)
		return
	}
	if buf[0] == 0x05 {
		ProxySocks5(conn)
	} else {
		ProxyHTTP(conn)
	}
}

func ProxyHTTP(conn net.Conn) {
	br := bufio.NewReader(conn)
	req, err := http.ReadRequest(br)
	if err != nil {
		log.Println("read request failed: ", err)
		return
	}
	defer req.Body.Close()

	if req.Method == http.MethodConnect {
		log.Println("http.MethodConnect: ", req.URL.Host)
		if err := handleHTTPConnect(conn, req); err != nil {
			log.Println("connect failed: ", err)
		}
		return
	}

	isHttpProxy := req.URL.IsAbs() && strings.HasPrefix(req.URL.Scheme, "http")
	if isHttpProxy {
		log.Println("isHttpProxy: ", req.URL.Host)
		if err := handleHTTPProxyForward(conn, req); err != nil {
			log.Println("http proxy forward failed: ", err)
		}
		return
	}
}

// handleHTTPConnect 处理 CONNECT（常见于 HTTPS），建立双向隧道。
func handleHTTPConnect(client net.Conn, req *http.Request) error {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	if host == "" {
		fmt.Fprintf(client, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return fmt.Errorf("missing host")
	}
	handle := transport.Local()
	if _, err := io.WriteString(client, "HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
		return err
	}
	handle(client, "tcp", host)
	return nil
}

// handleHTTPProxyForward 处理带绝对 URL 的 HTTP 代理请求（如 GET http://host/path）。
func handleHTTPProxyForward(client net.Conn, req *http.Request) error {
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}
	addr := net.JoinHostPort(req.URL.Hostname(), port)

	handle := transport.Local()
	remote, local := net.Pipe()
	go func() {
		handle(remote, "tcp", addr)
	}()

	// 复写请求
	out := req.Clone(req.Context())
	out.RequestURI = ""
	out.URL = &url.URL{
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}
	if out.URL.Path == "" {
		out.URL.Path = "/"
	}
	out.Host = req.URL.Host
	if err := out.Write(local); err != nil {
		return err
	}
	util.Pipe(client, local)
	return nil
}

func ProxySocks5(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read failed: ", err)
		return
	}
	log.Println(string(buf[:n]))
}
