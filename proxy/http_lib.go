package proxy

import (
	"fmt"
	"net"
)

// HTTP 代理头
var HttpProxyHeaders = []string{
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

// 单连接 Listener
type HttpSingleConnListener struct {
	conn net.Conn
}

func (l *HttpSingleConnListener) Accept() (net.Conn, error) {
	if l.conn == nil {
		return nil, fmt.Errorf("closed")
	}
	c := l.conn
	l.conn = nil
	return c, nil
}

func (l *HttpSingleConnListener) Close() error {
	return nil
}

func (l *HttpSingleConnListener) Addr() net.Addr {
	return l.conn.LocalAddr()
}
